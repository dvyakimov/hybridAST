package modules

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"hybridAST/core"
	"log"
	"net/http"
	"strings"
	"time"
)

func CheckArachni(url string) bool {
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

func SendStartArachni(host string) string {
	client := resty.New()
	respPostStartScan, err := client.R().
		SetBody(map[string]string{
			"url": host,
		}).
		SetHeader("Accept", "application/json").
		Post("http://arachni:7331/scans")

	if err != nil {
		log.Fatalf("ERROR:", err)
	}

	fmt.Println(string(respPostStartScan.Body()))

	lastId := gjson.Get(string(respPostStartScan.Body()), "id")

	for {
		respGetStatus, err := client.R().
			SetHeader("Accept", "application/json").
			Get("http://arachni:7331/scans/" + lastId.String())
		if err != nil {
			log.Fatalf("ERROR:", err)
		}
		lastMessage := gjson.Get(string(respGetStatus.Body()), "messages")
		lastStatus := gjson.Get(string(respGetStatus.Body()), "status")
		lastIssues := gjson.Get(string(respGetStatus.Body()), "issues")

		fmt.Println(lastStatus, lastMessage, lastIssues)
		if strings.Compare(lastStatus.String(), "done") == 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	respGetReport, err := client.R().
		SetHeader("Accept", "application/json").
		Get("http://arachni:7331/scans/" + lastId.String() + "/report.json")
	if err != nil {
		log.Fatalf("ERROR:", err)
	}

	fmt.Println(respGetReport.String())

	return respGetReport.String()
}

func StartScanArachni(url string) {
	AnalyzeArachni(SendStartArachni(url))
}

func ImportReportArachni(report string) {
	AnalyzeArachni(core.ImportReport(report))
}

func AnalyzeArachni(report string) {

	var db = core.InitDB()

	result := gjson.Get(report, "issues")

	for _, name := range result.Array() {

		var BugUrl = gjson.Get(name.String(), "page.dom.url").String()
		entry := core.Entrypoint{
			BugName: core.FindCWE(db, gjson.Get(name.String(), "name").String(),
				gjson.Get(name.String(), "cwe").String()),
			BugCWE:      gjson.Get(name.String(), "cwe").String(),
			BugHostPort: core.UrlExtractHostPort(BugUrl),
		}
		if core.UrlExtractPath(BugUrl) != "" {
			entry.BugPath = core.UrlExtractPath(BugUrl)
		} else {
			entry.BugPath = "/"
		}

		source := core.SourceData{
			Source:       "Arachni",
			Severity:     gjson.Get(name.String(), "severity").String(),
			SourceName:   gjson.Get(name.String(), "name").String(),
			SourceNameID: entry.ID,
		}
		core.UpdateEntry(entry, source, db, BugUrl)
	}
}
