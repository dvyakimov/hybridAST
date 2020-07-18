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

type rootNode struct {
	node  core.Node
	index int
}

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

func GraphIndexRoot(s []rootNode, e string) int {
	for i := 0; i < len(s); i++ {
		if s[i].node.String() == e {
			return s[i].index
		}
	}
	return -1
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	fileRead, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	content := string(fileRead)
	mapAs := make(map[string]string)

	r := *regexp.MustCompile(`(?m)^\s*from\s*(.*)\simport\s+(.*)\b`)
	rAs := *regexp.MustCompile(`(.*)\sas\s+(.*)`)
	//var rPath regexp.Regexp
	var arrayRoot []rootNode

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
				for j := 0; j < len(split)-1; j++ {
					var indexNode = core.FindInGraphByIndex(&g, indexRoot, split[j+1])
					if indexNode == -1 {
						g.AddNode(&core.Node{split[j+1]})
						g.AddEdge(&core.Node{split[j]}, &core.Node{split[j+1]})
					} else {
						indexRoot = indexNode
					}
				}
				/*-----------*/
				var regexpPathString string
				if mapAs[split[len(split)-1]] != "" {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + mapAs[importSplitAs[0][1]] + `\.(.+?)\,`
				} else {
					regexpPathString = `(?m)([[:word:]]*/|)\'\,\s*` + split[len(split)-1] + `\.(.+?)\,`
				}
				var rPath = *regexp.MustCompile(regexpPathString)
				//fmt.Println(rPath.String())
				var findPath = rPath.FindAllStringSubmatch(content, -1)
				for j := range findPath {
					//mt.Println("0",findPath[j][0])
					fmt.Println("1", findPath[j][1])
					//fmt.Println("2",findPath[j][2])

					var rFunc = *regexp.MustCompile(`(.*)\(.*\)|(.*)\)\s?|(.*)\s?`)
					var findRes = rFunc.FindAllStringSubmatch(findPath[j][2], -1)

					for l := 1; l < len(findRes[0]); l++ {
						if findRes[0][l] != "" {
							//fmt.Println(findRes[0][l])
							var splitFindRes = strings.Split(findRes[0][l], ".")
							g.AddNode(&core.Node{splitFindRes[0]})
							g.AddEdge(&core.Node{split[len(split)-1]}, &core.Node{splitFindRes[0]})
							//fmt.Println(splitFindRes[0])
							for m := 1; m <= len(splitFindRes)-1; m++ {
								//fmt.Println(splitFindRes[m])
								g.AddNode(&core.Node{splitFindRes[m]})
								g.AddEdge(&core.Node{splitFindRes[m-1]}, &core.Node{splitFindRes[m]})
							}
							//fmt.Println("splitFindRes[len(splitFindRes)-1]",splitFindRes[len(splitFindRes)-1])
							//g.AddNode(&core.Node{findPath[j][1]})
							//g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1]}, &core.Node{findPath[j][1]})
						}
					}

				}
				/*-----------*/
			} else {
				g.AddNode(&core.Node{split[0]})
				arrayRoot = append(arrayRoot, rootNode{core.Node{split[0]}, core.LastNode(&g)})
				//fmt.Println(arrayRoot)
				for j := 0; j < len(split)-1; j++ {
					g.AddNode(&core.Node{split[j+1]})
					g.AddEdge(&core.Node{split[j]}, &core.Node{split[j+1]})
				}
				/*-----------*/
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
					//fmt.Println(findPath[j][2])
					//fmt.Println("0",findPath[j][0])
					fmt.Println("1", findPath[j][1])
					//fmt.Println("2",findPath[j][2])
					var rFunc = *regexp.MustCompile(`(.*)\(.*\)|(.*)\)\s?|(.*)\s?`)
					var findRes = rFunc.FindAllStringSubmatch(findPath[j][2], -1)

					for l := 1; l < len(findRes[0]); l++ {
						if findRes[0][l] != "" {
							//fmt.Println(findRes[0][l])
							var splitFindRes = strings.Split(findRes[0][l], ".")
							g.AddNode(&core.Node{splitFindRes[0]})
							g.AddEdge(&core.Node{split[len(split)-1]}, &core.Node{splitFindRes[0]})
							//fmt.Println(splitFindRes[0])
							for m := 1; m <= len(splitFindRes)-1; m++ {
								//fmt.Println(splitFindRes[m])
								g.AddNode(&core.Node{splitFindRes[m]})
								g.AddEdge(&core.Node{splitFindRes[m-1]}, &core.Node{splitFindRes[m]})
							}
							//fmt.Println("splitFindRes[len(splitFindRes)-1]",splitFindRes[len(splitFindRes)-1])
							//g.AddNode(&core.Node{findPath[j][1]})
							//g.AddEdge(&core.Node{splitFindRes[len(splitFindRes)-1]}, &core.Node{findPath[j][1]})
						}
					}

					//fmt.Println(findPath)
				}
				/*-----------*/

			}

			g.String()
			//fmt.Println("---------")

		}

	}

	return lines, scanner.Err()
}

var g core.ItemGraph

func StartAnalyzeBandit() {
	var db = core.InitDB()

	result := gjson.Get(core.ImportReport("examples-report/bandit-django.json"), "results")

	/*-----------testing parsing ------------*/
	lines, err := readLines("tempSource/urls.py")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	for i, line := range lines {
		fmt.Println(i, line)
	}

	/*----------*/

	for _, name := range result.Array() {
		var CweId = TestidToCWE(core.GetValue(name, "test_id"))

		entry := core.Entrypoint{
			BugName: core.FindCWE(db, gjson.Get(name.String(), "issue_text").String(), CweId),
			BugCWE:  CweId,
		}
		/*--------------------*/

		/*-------------------*/
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
