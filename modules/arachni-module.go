package modules

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"hybridAST/entrypoint"
	"log"
	"strings"
	"time"
)

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

	db, err := gorm.Open("mysql", "root:root@/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&entrypoint.Entrypoint{}, &entrypoint.Params{}, &entrypoint.SourceData{})
	db.Model(&entrypoint.Params{}).AddForeignKey("params_id", "entrypoints(id)", "CASCADE", "CASCADE")
	db.Model(&entrypoint.SourceData{}).AddForeignKey("source_name_id", "entrypoints(id)", "CASCADE", "CASCADE")

	var CweResult []*entrypoint.CweList

	result := gjson.Get(entrypoint.ImportReport("examples-report/arachni-report-example.json"), "issues")
	//result := gjson.Get(StartScanArachni("http://127.0.0.1:8000"), "issues")

	for _, name := range result.Array() {

		var BugUrl = gjson.Get(name.String(), "page.dom.url").String()
		entry := entrypoint.Entrypoint{
			BugName: entrypoint.FindCWE(db, gjson.Get(name.String(), "name").String(),
				gjson.Get(name.String(), "cwe").String()),
			BugCWE:      gjson.Get(name.String(), "cwe").String(),
			BugHostPort: entrypoint.UrlExtractHostPort(BugUrl),
		}
		if entrypoint.UrlExtractPath(BugUrl) != "" {
			entry.BugPath = entrypoint.UrlExtractPath(BugUrl)
		} else {
			entry.BugPath = "/"
		}

		if len(CweResult) > 0 {
			entry.BugName = CweResult[0].Name
		} else {
			entry.BugName = gjson.Get(name.String(), "name").String()
		}

		source := entrypoint.SourceData{
			Source:       "Arachni",
			Severity:     gjson.Get(name.String(), "severity").String(),
			SourceName:   gjson.Get(name.String(), "name").String(),
			SourceNameID: entry.ID,
		}
		entrypoint.UpdateEntry(entry, source, db, BugUrl)
	}
}
