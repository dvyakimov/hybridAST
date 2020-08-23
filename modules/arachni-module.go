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
	defer resp.Body.Close()

	if resp != nil {
		return true
	} else {
		return false
	}
}

func StartScanArachni(host string) string {
	client := resty.New()
	respPostStartScan, err := client.R().
		SetBody(map[string]string{
			"url": host,
		}).
		SetHeader("Accept", "application/json").
		Post("http://192.168.168.3:7331/scans")

	if err != nil {
		log.Fatalf("ERROR:", err)
	}

	lastId := gjson.Get(string(respPostStartScan.Body()), "id")

	for {
		respGetStatus, err := client.R().
			SetHeader("Accept", "application/json").
			Get("http://192.168.168.3:7331/scans/" + lastId.String())
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
		Get("http://192.168.168.3:7331/scans/" + lastId.String() + "/report.json")
	if err != nil {
		log.Fatalf("ERROR:", err)
	}
	return respGetReport.String()
}

func StartAnalyzeArachni() {
	var db = core.InitDB()

	result := gjson.Get(core.ImportReport("examples-report/arachni-report-example.json"), "issues")
	//result := gjson.Get(StartScanArachni("http://127.0.0.1:8000"), "issues")

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
