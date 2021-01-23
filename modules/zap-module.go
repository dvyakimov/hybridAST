package modules

import (
	"encoding/xml"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zaproxy/zap-api-go/zap"
	"hybridAST/core"
	//"io"
	"log"
	"net/http"
	"strconv"
	//"strings"
	"time"
	//"github.com/antchfx/xmlquery"
)

func CheckZAP(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
	}
	//defer resp.Body.Close()

	if resp != nil {
		return true
	} else {
		return false
	}
}

func SendStartZap(host string) string {
	target := host
	cfg := &zap.Config{
		Proxy: "http://zaproxy:8090",
	}
	client, err := zap.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Start spidering the target
	fmt.Println("Spider : " + target)
	resp, err := client.Spider().Scan(target, "", "", "", "")
	if err != nil {
		log.Fatal(err)
	}

	// The scan now returns a scan id to support concurrent scanning
	scanid := resp["scan"].(string)
	for {
		time.Sleep(1000 * time.Millisecond)
		resp, _ = client.Spider().Status(scanid)
		progress, _ := strconv.Atoi(resp["status"].(string))
		if progress >= 100 {
			break
		}
	}
	fmt.Println("Spider complete")

	// Give the passive scanner a chance to complete
	time.Sleep(2000 * time.Millisecond)

	fmt.Println("Active scan : " + target)
	resp, err = client.Ascan().Scan(target, "True", "False", "", "", "", "")
	if err != nil {
		log.Fatal(err)
	}
	// The scan now returns a scan id to support concurrent scanning
	scanid = resp["scan"].(string)
	for {
		time.Sleep(5000 * time.Millisecond)
		resp, _ = client.Ascan().Status(scanid)
		progress, _ := strconv.Atoi(resp["status"].(string))
		fmt.Printf("Active Scan progress : %d\n", progress)
		if progress >= 100 {
			break
		}
	}
	fmt.Println("Active Scan complete")
	RawReport, err := client.Core().Jsonreport()
	if err != nil {
		log.Fatal(err)
	}

	return string(RawReport)
}

func StartScanZap(url string, appId uint, SeverityFlag bool) {
	fmt.Println("Start Scan is completed")
	SendStartScanResult := SendStartZap(url)
	if SendStartScanResult != "" {
		fmt.Println("Send Start ZAP is completed")
		AnalyzeZapJson(SendStartScanResult, appId, SeverityFlag)
	} else {
		return
	}
}

func ImportReportZapJson(report string, appId uint, SeverityFlag bool) {
	AnalyzeZapJson(core.ImportReport(report), appId, SeverityFlag)
}

func ImportReportZapXml(report string, appId uint, SeverityFlag bool) {
	AnalyzeZapXml(core.ImportReport(report), appId, SeverityFlag)
}

// array of all Users in the file
type OWASPReport struct {
	XMLName xml.Name `xml:"OWASPZAPReport"`
	Site    []Site   `xml:"site"`
}

type Site struct {
	XMLName xml.Name `xml:"site"`
	Alerts  []Alerts `xml:"alerts"`
}

type Alerts struct {
	XMLName   xml.Name    `xml:"alerts"`
	Alertitem []Alertitem `xml:"alertitem"`
}
type Alertitem struct {
	XMLName    xml.Name `xml:"alertitem"`
	Pluginid   string   `xml:"pluginid"`
	Alert      string   `xml:"alert"`
	Riskcode   string   `xml:"riskcode"`
	Confidence string   `xml:"confidence"`
	Riskdesc   string   `xml:"riskdesc"`
	Desc       string   `xml:"desc"`
	Uri        string   `xml:"uri"`
	Param      string   `xml:"param"`
	Attack     string   `xml:"attack"`
	Otherinfo  string   `xml:"otherinfo"`
	Solution   string   `xml:"solution"`
	Evidence   string   `xml:"evidence"`
	Reference  string   `xml:"reference"`
	Cweid      string   `xml:"cweid"`
}

func AnalyzeZapXml(report string, appId uint, SeverityFlag bool) {
	//fmt.Println(report)
	var db = core.InitDB()
	var structedReport OWASPReport

	err := xml.Unmarshal([]byte(report), &structedReport)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	for _, name := range structedReport.Site[0].Alerts[0].Alertitem {

		if SeverityFlag == true && name.Riskcode != "3" {
			continue
		}

		fmt.Printf("Value: %v\n", name.Alert)

		entry := core.Entrypoint{
			BugName:     core.FindNameByCWE(db, name.Alert, name.Cweid),
			BugCWE:      name.Cweid,
			BugHostPort: core.UrlExtractHostPort(name.Uri),
			AppId:       appId,
		}

		if core.UrlExtractPath(name.Uri) != "" {
			entry.BugPath = core.UrlExtractPath(name.Uri)
		} else {
			entry.BugPath = "/"
		}

		source := core.SourceData{
			Source:       "OWASP ZAP",
			Severity:     name.Riskcode,
			SourceName:   name.Alert,
			SourceNameID: entry.ID,
		}

		params := core.UrlExtractParametr(name.Uri)

		core.UpdateEntry(entry, source, db, params) //Убрать из аргументов BugURL

	}

}

func AnalyzeZapJson(report string, appId uint, SeverityFlag bool) {
	fmt.Println(report)

	var db = core.InitDB()

	result := gjson.Get(report, "site.0.alerts")

	for _, name := range result.Array() {
		resultInstances := gjson.Get(name.String(), "instances")
		for _, nameSecond := range resultInstances.Array() {
			if SeverityFlag == true && gjson.Get(name.String(), "riskcode").String() != "3" {
				continue
			}

			var BugUrl = gjson.Get(nameSecond.String(), "uri").String()
			entry := core.Entrypoint{
				BugName: core.FindNameByCWE(db, gjson.Get(name.String(), "name").String(),
					gjson.Get(name.String(), "cweid").String()),
				BugCWE:      gjson.Get(name.String(), "cweid").String(),
				BugHostPort: core.UrlExtractHostPort(BugUrl),
				AppId:       appId,
			}

			if core.UrlExtractPath(BugUrl) != "" {
				entry.BugPath = core.UrlExtractPath(BugUrl)
			} else {
				entry.BugPath = "/"
			}

			source := core.SourceData{
				Source:       "OWASP ZAP",
				Severity:     gjson.Get(name.String(), "riskcode").String(),
				SourceName:   gjson.Get(name.String(), "name").String(),
				SourceNameID: entry.ID,
			}

			params := core.UrlExtractParametr(BugUrl)

			core.UpdateEntry(entry, source, db, params) //Убрать из аргументов BugURL
		}
	}
}
