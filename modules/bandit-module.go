package modules

import (
	"fmt"
	"github.com/tidwall/gjson"
	"hybridAST/core"
)

func TestidToCWE(testid string) string {
	lines, err := core.ReadCsv("cwe/bandit-cwe.csv")

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
	var db = core.InitDB()

	result := gjson.Get(core.ImportReport("examples-report/bandit-django.json"), "results")

	for _, name := range result.Array() {
		var CweId = TestidToCWE(core.GetValue(name, "test_id"))

		entry := core.Entrypoint{
			BugName: core.FindCWE(db, gjson.Get(name.String(), "issue_text").String(), CweId),
			BugCWE:  CweId,
		}

		source := core.SourceData{
			Source:       "Bandit",
			Severity:     gjson.Get(name.String(), "issue_severity").String(),
			SourceName:   gjson.Get(name.String(), "issue_text").String(),
			SourceNameID: entry.ID,
			CodeLine:     gjson.Get(name.String(), "line_number").String(),
			CodeFile:     gjson.Get(name.String(), "filename").String(),
		}
		core.UpdateEntry(entry, source, db, "")

		fmt.Println(entry)
		fmt.Println(source)
	}
}
