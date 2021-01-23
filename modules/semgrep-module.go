package modules

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"hybridAST/core"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PathSast struct {
	PathValue string
	FileValue string
	FuncValue string
	StartLine int
	EndLine   int
}

type VulnSast struct {
	BugName     string
	BugSeverity string
	BugCWE      string

	FileValue string
	StartLine int
	PathValue string
}

func GetSemgrepCWE(name string) string {
	var re = regexp.MustCompile(`(?m)CWE-(\d.*):`)

	res := re.FindAllStringSubmatch(name, -1)

	return res[0][1]
}

// TODO: помещать значения в базу данных
func SemgrepGetPath(filename string) map[string]PathSast {
	out, err := exec.Command("semgrep", "--config", "./internal-rules/java-spring-router.yml", "-v", "--json", "-o", "report.json", "./output-folder").Output()

	if err != nil {
		log.Fatal(err)
	}

	result := gjson.Get(string(out), "results")

	pathMap := make(map[string]PathSast)

	for _, name := range result.Array() {
		startline, err := strconv.Atoi(gjson.Get(name.String(), "start.line").String())
		if err != nil {
			log.Fatal(err)
		}
		endline, err := strconv.Atoi(gjson.Get(name.String(), "end.line").String())
		if err != nil {
			log.Fatal(err)
		}

		pathTemp := PathSast{
			PathValue: gjson.Get(name.String(), "extra.metavars.$PATH.abstract_content").String(),
			FileValue: gjson.Get(name.String(), "path").String(),
			FuncValue: gjson.Get(name.String(), "extra.metavars.$FUNC.abstract_content").String(),
			StartLine: startline,
			EndLine:   endline,
		}

		pathMap[gjson.Get(name.String(), "path").String()] = pathTemp

	}

	return pathMap
}

func GetContextRoot(db *gorm.DB, appid uint) string {

	var apptemp core.AppList

	db.Find(&apptemp, "id=?", appid)

	return apptemp.ContextRoot

}

func SemgrepGetBugResult() string {
	//TODO: сделать возможность вставлять правила из формы и проверить соединение с Интернет
	out, err := exec.Command("semgrep", "--config", "p/findsecbugs", "--json", "-o", "result-scan.json", "./output-folder").Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func SemgrepScan(filename string, AppId uint, SeverityFlag bool) {

	var db = core.InitDB()

	pathMap := SemgrepGetPath(filename) // TODO: забирать значения из базы данных

	fmt.Println(pathMap)

	BugResult := SemgrepGetBugResult() //TODO: сделать возмоджность получения отчет через загрузку без запуска сканирования

	ContextRootValue := GetContextRoot(db, AppId) // TODO: сделать поиск context-root через foreignkey таблицы

	var VulnSastArray []VulnSast

	resultVulnt := gjson.Get(BugResult, "results")

	for _, name := range resultVulnt.Array() {
		if SeverityFlag == true && gjson.Get(name.String(), "extra.severity").String() != "WARNING" {
			continue
		}
		startline, err := strconv.Atoi(gjson.Get(name.String(), "start.line").String())
		if err != nil {
			log.Fatal(err)
		}
		fileValue := gjson.Get(name.String(), "path").String()

		bugValue := gjson.Get(name.String(), "extra.metadata.cwe").String()

		//TODO: VulnSast можно будет потом выпилить
		vulnTemp := VulnSast{
			BugName:     core.FindNameByCWE(db, bugValue, GetSemgrepCWE(bugValue)),
			BugCWE:      GetSemgrepCWE(bugValue),
			BugSeverity: gjson.Get(name.String(), "extra.severity").String(),
			FileValue:   gjson.Get(name.String(), "path").String(),
			StartLine:   startline,
		}

		entry := core.Entrypoint{
			BugName: core.FindNameByCWE(db, bugValue, GetSemgrepCWE(bugValue)),
			BugCWE:  GetSemgrepCWE(bugValue),
			AppId:   AppId,
		}

		if startline <= pathMap[fileValue].EndLine && startline >= pathMap[fileValue].StartLine {
			entry.BugPath = "/" + ContextRootValue + strings.ReplaceAll(pathMap[fileValue].PathValue, "\"", ``)
		} // TODO: else рекурсивно искать к какому хендлеру подходит эта функция

		source := core.SourceData{
			Source:       "Semgrep",
			Severity:     gjson.Get(name.String(), "extra.severity").String(),
			SourceName:   core.FindNameByCWE(db, bugValue, GetSemgrepCWE(bugValue)),
			SourceNameID: entry.ID,
			CodeLine:     strconv.Itoa(startline),
			CodeFile:     gjson.Get(name.String(), "path").String(),
		}

		core.UpdateEntry(entry, source, db, nil)

		VulnSastArray = append(VulnSastArray, vulnTemp)
	}

	fmt.Println(VulnSastArray)

	DeleteReportFile(filename)

}

func DeleteReportFile(filename string) {
	_, err := exec.Command("rm", "-rf", "./output-folder").Output()

	if err != nil {
		log.Fatal(err)
	}
}
