package modules

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"hybridAST/entrypoint"
	"os"
)

func ReadCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func TestidToCWE(testid string) string {
	lines, err := ReadCsv("cwe/bandit-cwe.csv")

	if err != nil {
		fmt.Println(err)
	}
	for _, line := range lines {
		if line[0] == testid {
			return line[1]
		}
	}
	return "0"
}

func StartAnalyzeBandit() {
	db, err := gorm.Open("mysql", "root:root@/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&entrypoint.Entrypoint{}, &entrypoint.Params{}, &entrypoint.SourceData{})
	db.Model(&entrypoint.Params{}).AddForeignKey("params_id", "entrypoints(id)", "CASCADE", "CASCADE")
	db.Model(&entrypoint.SourceData{}).AddForeignKey("source_name_id", "entrypoints(id)", "CASCADE", "CASCADE")

	result := gjson.Get(entrypoint.ImportReport("examples-report/bandit-django.json"), "results")

	for _, name := range result.Array() {
		var CweId = TestidToCWE(gjson.Get(name.String(), "test_id").String())

		entry := entrypoint.Entrypoint{
			BugName: entrypoint.FindCWE(db, gjson.Get(name.String(), "issue_text").String(), CweId),
			BugCWE:  gjson.Get(name.String(), "cwe").String(),
		}

		source := entrypoint.SourceData{
			Source:       "Bandit",
			Severity:     gjson.Get(name.String(), "issue_severity").String(),
			SourceName:   gjson.Get(name.String(), "issue_text").String(),
			SourceNameID: entry.ID,
			CodeLine:     gjson.Get(name.String(), "line_number").String(),
			CodeFile:     gjson.Get(name.String(), "filename").String(),
		}
		entrypoint.UpdateEntry(entry, source, db, "")

		fmt.Println(entry)
		fmt.Println(source)
	}
}
