package modules

import (
	"bufio"
	"fmt"
	"github.com/tidwall/gjson"
	"hybridAST/core"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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

func GraphIndexRoot(s []core.Node, e string) int {
	for i := 0; i < len(s); i++ {
		if s[i].String() == e {
			return s[i].ID
		}
	}
	return -1
}

func buildGraph(path string) ([]core.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileRead, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	content := string(fileRead)
	mapAs := make(map[string]string)

	r := *regexp.MustCompile(`(?m)^\s*from\s*(.*)\simport\s+(.*)\b`)
	rAs := *regexp.MustCompile(`(.*)\sas\s+(.*)`)

	var arrayRoot []core.Node
	var NodeId = -1

	for scanner.Scan() {

		res := r.FindAllStringSubmatch(scanner.Text(), -1)

		for i := range res {

			var importSplitAs = rAs.FindAllStringSubmatch(res[i][2], -1)

			var split = strings.Split(res[i][1], ".")
			if importSplitAs != nil {
				split = append(split, importSplitAs[0][1])
				mapAs[importSplitAs[0][1]] = importSplitAs[0][2]
			} else {
				split = append(split, res[i][2])
			}

			var indexRoot = GraphIndexRoot(arrayRoot, split[0])
			if indexRoot != -1 {
				var LastNodeId int
				for j := 0; j < len(split)-1; j++ {
					var indexNode = core.FindInGraphByIndex(&g, indexRoot, split[j+1])
					if indexNode == -1 {
						NodeId++
						g.AddNode(&core.Node{split[j+1], false, NodeId})
						g.AddEdge(&core.Node{split[j], false, indexRoot}, &core.Node{split[j+1], false, NodeId})
						indexRoot = NodeId
					} else {
						indexRoot = indexNode
					}
				}
				LastNodeId = NodeId
				var regexpPathString string
				if mapAs[split[len(split)-1]] != "" {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + mapAs[importSplitAs[0][1]] + `\.(.+?)\,`
				} else {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + split[len(split)-1] + `\.(.+?)\,`
				}
				var rPath = *regexp.MustCompile(regexpPathString)
				var findPath = rPath.FindAllStringSubmatch(content, -1)
				for j := range findPath {

					var rFunc = *regexp.MustCompile(`(.*)\(.*\)|(.*)\)\s?|(.*)\s?`)
					var findRes = rFunc.FindAllStringSubmatch(findPath[j][2], -1)

					for l := 1; l < len(findRes[0]); l++ {
						if findRes[0][l] != "" {
							var splitFindRes = strings.Split(findRes[0][l], ".")
							NodeId++
							g.AddNode(&core.Node{splitFindRes[0], true, NodeId})
							g.AddEdge(&core.Node{split[len(split)-1], false, LastNodeId}, &core.Node{splitFindRes[0], true, NodeId})
							for m := 1; m <= len(splitFindRes)-1; m++ {
								NodeId++
								g.AddNode(&core.Node{splitFindRes[m], true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[m-1], true, NodeId - 1}, &core.Node{splitFindRes[m], true, NodeId})
							}
							if findPath[j][1] == "" {
								NodeId++
								g.AddNode(&core.Node{"/", true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1], true, NodeId - 1}, &core.Node{"/", true, NodeId})
							} else {
								NodeId++
								g.AddNode(&core.Node{findPath[j][1], true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1], true, NodeId - 1}, &core.Node{findPath[j][1], true, NodeId})
							}
						}
					}
				}
			} else {
				NodeId++
				g.AddNode(&core.Node{split[0], false, NodeId})
				arrayRoot = append(arrayRoot, core.Node{split[0], false, NodeId})
				var LastNodeId int
				for j := 0; j < len(split)-1; j++ {
					NodeId++
					g.AddNode(&core.Node{split[j+1], false, NodeId})
					g.AddEdge(&core.Node{split[j], false, NodeId - 1}, &core.Node{split[j+1], false, NodeId})
					LastNodeId = NodeId
				}
				var regexpPathString string
				if mapAs[split[len(split)-1]] != "" {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + mapAs[importSplitAs[0][1]] + `\.(.+?)\,`
				} else {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + split[len(split)-1] + `\.(.+?)\,`
				}
				var rPath = *regexp.MustCompile(regexpPathString)
				//fmt.Println(rPath.String())
				var findPath = rPath.FindAllStringSubmatch(content, -1) // тут не scanner.text
				for j := range findPath {
					var rFunc = *regexp.MustCompile(`(.*)\(.*\)|(.*)\)\s?|(.*)\s?`)
					var findRes = rFunc.FindAllStringSubmatch(findPath[j][2], -1)

					for l := 1; l < len(findRes[0]); l++ {
						if findRes[0][l] != "" {
							var splitFindRes = strings.Split(findRes[0][l], ".")
							NodeId++
							g.AddNode(&core.Node{splitFindRes[0], true, NodeId})
							g.AddEdge(&core.Node{split[len(split)-1], false, LastNodeId}, &core.Node{splitFindRes[0], true, NodeId})
							for m := 1; m <= len(splitFindRes)-1; m++ {
								NodeId++
								g.AddNode(&core.Node{splitFindRes[m], true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[m-1], true, NodeId - 1}, &core.Node{splitFindRes[m], true, NodeId})
							}
							if findPath[j][1] == "" {
								NodeId++
								g.AddNode(&core.Node{"/", true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1], true, NodeId - 1}, &core.Node{"/", true, NodeId})
							} else {
								NodeId++
								g.AddNode(&core.Node{findPath[j][1], true, NodeId})
								g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1], true, NodeId - 1}, &core.Node{findPath[j][1], true, NodeId})
							}
						}
					}
				}
			}
		}
	}
	return arrayRoot, err
}

func GetFuncName(path string, vulnline int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	var re = regexp.MustCompile(`(?m)def\s+(.*)\(.*\):.*`)
	var index int64
	var funcName string
	for scanner.Scan() {
		index++
		if re.FindString(scanner.Text()) != "" {
			var findRes = re.FindAllStringSubmatch(scanner.Text(), -1)
			funcName = findRes[0][1]
		}
		if index >= vulnline {
			return funcName, err
		}
	}
	return "", err
}

func findPathInGraph(g core.ItemGraph, funcname string, filename string, arrayRoot []core.Node) string {
	filename = strings.Replace(filename, ".py", "", -1)
	splitSymbol := strings.Split(filename, "/")

	for i := 0; i <= len(splitSymbol)-1; i++ {
		if splitSymbol[i] != "." {
			for j := 0; j < len(arrayRoot)-1; j++ {
				if arrayRoot[j].String() == splitSymbol[i] {
					var indexRoot = arrayRoot[j].ID
					var indexNode int
					var isFound = false
					for ; i < len(splitSymbol)-1; i++ {
						indexNode = core.FindInGraphByIndex(&g, indexRoot, splitSymbol[i+1])
						if indexNode != -1 {
							indexRoot = indexNode
							isFound = true
						} else {
							isFound = false
						}
					}
					if isFound == true {
						indexNode = core.FindInGraphByIndex(&g, indexRoot, funcname)
						if indexNode != -1 {
							if g.GetTheLastNodeValueString(indexNode) != "" {
								return g.GetTheLastNodeValueString(indexNode)
							} else {
								return ""
							}
						} else {
							return ""
						}
					}
				}
			}
		}
	}
	return ""
}

//Вторник:
// 1. Доделать path в графе
// 2. Импорт проекта и парсинг по файлам
// 3. Идеи для поиска параметра URL, сформированный план

//Среда:
// Реализация поиска параметра для URL для моего примера
// Запаковывание инструментов сканирования в Docker + запаковывание своего приложения

// Четверг:
// Реализация веб-интерфейса

// Пятница:
// Реализация графиков

var g core.ItemGraph

func GetUrlPath(pathProject string) (string, error) {
	managefile, err := os.Open(pathProject + "/manage.py")
	if err != nil {
		return "", err
	}
	defer managefile.Close()

	// добавить поиск в конструкции типа
	// var re = regexp.MustCompile(`os\.environ\[\'DJANGO_SETTINGS_MODULE\'\]\s?\=\s?\'\s?(.*)\'.*`)

	var re = regexp.MustCompile(`(?U)os\.environ\.setdefault\(\'DJANGO_SETTINGS_MODULE\',\s+\'(.*)\'.*`)
	scanner := bufio.NewScanner(managefile)
	var findResSetting [][]string
	for scanner.Scan() {
		if re.FindString(scanner.Text()) != "" {
			findResSetting = re.FindAllStringSubmatch(scanner.Text(), -1)
			//fmt.Println(findResSetting[0][1])
		}
	}
	var settingsString = pathProject
	if findResSetting != nil {
		splitSymbol := strings.Split(findResSetting[0][1], ".")
		for i := range splitSymbol {
			settingsString = settingsString + "/" + splitSymbol[i]
		}
		settingsString = settingsString + ".py"
	}
	//fmt.Println(settingsString)

	settingsfile, err := os.Open(settingsString)
	if err != nil {
		return "", err
	}
	defer settingsfile.Close()

	urlFinding := regexp.MustCompile(`(?mU)^ROOT_URLCONF\s?\=\s?\'(.*)\'`)
	scannerSettings := bufio.NewScanner(settingsfile)
	var findResUrl [][]string

	for scannerSettings.Scan() {
		if urlFinding.FindString(scannerSettings.Text()) != "" {
			findResUrl = urlFinding.FindAllStringSubmatch(scannerSettings.Text(), -1)
		}
	}

	var urlString = pathProject

	if findResUrl != nil {
		splitSymbol := strings.Split(findResUrl[0][1], ".")
		for i := range splitSymbol {
			urlString = urlString + "/" + splitSymbol[i]
		}
		urlString = urlString + ".py"
	} else {
		return "", err
	}

	return urlString, err
}

func GetParam(path string, funcname string, vulnNumber int64) (string, error) {

	var funcStart = regexp.MustCompile(`(?m)def\s+` + funcname + `\(.*\):`)
	var funcEnd = regexp.MustCompile(`(?m)return\s+.*`)
	var valueGet = regexp.MustCompile(`(?Um)\s?([[:word:]].*[[:word:]])\s\=\srequest\.GET.*\'\s?.*([[:word:]].*[[:word:]])\s?\'\)`)
	var ifvaluetrue *regexp.Regexp
	var elsevaluetrue = regexp.MustCompile(`(?m)else:`)

	sourcecodefile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer sourcecodefile.Close()

	scanner := bufio.NewScanner(sourcecodefile)
	//var funccode [][]string
	var valueString [][]string
	var isFuncBlock bool
	var isIfBlock bool
	var index int64
	for scanner.Scan() {
		index++
		if funcStart.FindString(scanner.Text()) != "" {
			isFuncBlock = true
		}
		if funcEnd.FindString(scanner.Text()) != "" {
			isFuncBlock = false
		}
		if valueGet.FindString(scanner.Text()) != "" && isFuncBlock == true {
			valueString = valueGet.FindAllStringSubmatch(scanner.Text(), -1)
			//fmt.Println(valueString[0][1])
			ifvaluetrue = regexp.MustCompile(`(?m)if\s+` + valueString[0][1] + `:`)
		}
		if ifvaluetrue != nil {
			if ifvaluetrue.FindString(scanner.Text()) != "" && isFuncBlock == true {
				isIfBlock = true
			}
		}
		if elsevaluetrue.FindString(scanner.Text()) != "" && isFuncBlock == true {
			isIfBlock = false
		}
		if isIfBlock == true && isFuncBlock == true && index == vulnNumber {
			return valueString[0][2], err
		}
	}
	return "", err
}

func StartAnalyzeBandit() {
	var db = core.InitDB()

	var urlFile, errurlFile = GetUrlPath("/Users/denisyakimov/Desktop/django-vuln-example")

	if errurlFile != nil {
		log.Fatalf("urlFile: %s", errurlFile)
	}

	result := gjson.Get(core.ImportReport("examples-report/bandit-django.json"), "results")

	var arrayRoot, err = buildGraph(urlFile)
	if err != nil {
		log.Fatalf("buildGraph: %s", err)
	}

	/*----------*/

	for _, name := range result.Array() {
		var BugUrl = "http://127.0.0.1"
		var CweId = TestidToCWE(core.GetValue(name, "test_id"))

		entry := core.Entrypoint{
			BugName: core.FindNameByCWE(db, gjson.Get(name.String(), "issue_text").String(), CweId),
			BugCWE:  CweId,
		}

		/*-----------testing parsing ------------*/

		var funcresult string
		funcresult, err = GetFuncName(gjson.Get(name.String(), "filename").String(), gjson.Get(name.String(), "line_number").Int())
		if err != nil {
			log.Fatalf("GetFuncName: %s", err)
		}

		var BugPath = findPathInGraph(g, funcresult, gjson.Get(name.String(), "filename").String(), arrayRoot)

		if BugPath != "" {
			entry.BugPath = BugPath
		}
		/*--------------------*/
		var BugParam string
		BugParam, err = GetParam(gjson.Get(name.String(), "filename").String(), funcresult, gjson.Get(name.String(), "line_number").Int())

		if BugParam != "" {
			BugUrl = BugUrl + BugPath + "?" + BugParam + "="
		} else {
			BugUrl = BugUrl + BugPath
		}

		//fmt.Println(BugParam)

		fmt.Println(BugUrl)

		source := core.SourceData{
			Source:       "Bandit",
			Severity:     gjson.Get(name.String(), "issue_severity").String(),
			SourceName:   gjson.Get(name.String(), "issue_text").String(),
			SourceNameID: entry.ID,
			CodeLine:     gjson.Get(name.String(), "line_number").String(),
			CodeFile:     gjson.Get(name.String(), "filename").String(),
		}
		core.UpdateEntry(entry, source, db, nil)

		fmt.Println(entry)
		fmt.Println(source)
	}
}
