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
		g.String()
	}
	return arrayRoot, err
}

func readLines(path string, vulnline int64) (string, error) {
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

func findPathInGraph(g core.ItemGraph, funcname string, filename string, arrayRoot []core.Node) {
	splitSymbol := strings.Split(filename, "/")
	for i := range splitSymbol {
		if splitSymbol[i] != "." {
			for j := 0; j < len(arrayRoot)-1; j++ {
				if arrayRoot[j].String() == splitSymbol[i] {
					var indexRoot = arrayRoot[j].ID

					//fmt.Println(g.FindInGraph(splitSymbol[i+1]))

					var indexNode = core.FindInGraphByIndex(&g, indexRoot, splitSymbol[i+1])
					fmt.Println(indexRoot)
					fmt.Println(splitSymbol[i+1])

					if indexNode != -1 {
						fmt.Println(indexNode)
					} else {
						//indexRoot = indexNode
					}
				}
			}
		}
	}
}

var g core.ItemGraph

func StartAnalyzeBandit() {
	var db = core.InitDB()

	result := gjson.Get(core.ImportReport("examples-report/bandit-django.json"), "results")

	var arrayRoot, err = buildGraph("tempSource/urls.py")
	if err != nil {
		log.Fatalf("buildGraph: %s", err)
	}
	fmt.Println(arrayRoot)

	/*----------*/

	for _, name := range result.Array() {
		var CweId = TestidToCWE(core.GetValue(name, "test_id"))

		entry := core.Entrypoint{
			BugName: core.FindCWE(db, gjson.Get(name.String(), "issue_text").String(), CweId),
			BugCWE:  CweId,
		}
		/*--------------------*/

		/*-----------testing parsing ------------*/

		var funcresult string
		funcresult, err = readLines("tempSource/blog_views.py", gjson.Get(name.String(), "line_number").Int())
		if err != nil {
			log.Fatalf("readLines: %s", err)
		}
		fmt.Println(funcresult)

		fmt.Println(gjson.Get(name.String(), "filename").String())

		findPathInGraph(g, funcresult, gjson.Get(name.String(), "filename").String(), arrayRoot)

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
