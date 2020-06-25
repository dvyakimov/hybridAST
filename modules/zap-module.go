package modules

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zaproxy/zap-api-go/zap"
	"hybridAST/core"
	"log"
	"strconv"
	"time"
)

var target string

func init() {
	flag.StringVar(&target, "target", "http://192.168.168.2:8080", "target address")
	flag.Parse()
}

func StartScanZAP() string {
	cfg := &zap.Config{
		Proxy: "http://192.168.168.3:8080",
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

func StartAnalyzeZAP() {
	var db = core.InitDB()

	result := gjson.Get(core.ImportReport("examples-report/zap-report-example.json"), "site.0.alerts")

	for _, name := range result.Array() {
		resultInstances := gjson.Get(name.String(), "instances")
		for _, nameSecond := range resultInstances.Array() {

			var BugUrl = gjson.Get(nameSecond.String(), "uri").String()
			entry := core.Entrypoint{
				BugName: core.FindCWE(db, gjson.Get(name.String(), "name").String(),
					gjson.Get(name.String(), "cweid").String()),
				BugCWE:      gjson.Get(name.String(), "cweid").String(),
				BugHostPort: core.UrlExtractHostPort(BugUrl),
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
			core.UpdateEntry(entry, source, db, BugUrl) //Убрать из аргументов BugURL
		}
	}
}
