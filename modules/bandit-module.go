package modules

import (
	"bufio"
	"fmt"
	"github.com/tidwall/gjson"
	"hybridAST/core"
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
	scanner := bufio.NewScanner(file)
	r := *regexp.MustCompile(`(?m)^\s*from\s*(.*)\simport\s+(.*)\b`) // this can also be a regex
	var arrayRoot []rootNode
	for scanner.Scan() {

		res := r.FindAllStringSubmatch(scanner.Text(), -1)

		for i := range res {
			fmt.Printf("first: %s, second: %s\n", res[i][1], res[i][2])
			var split = strings.Split(res[i][1], ".")
			split = append(split, res[i][2])
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
			} else {
				g.AddNode(&core.Node{split[0]})
				arrayRoot = append(arrayRoot, rootNode{core.Node{split[0]}, core.LastNode(&g)})
				fmt.Println(arrayRoot)
				for j := 0; j < len(split)-1; j++ {
					g.AddNode(&core.Node{split[j+1]})
					g.AddEdge(&core.Node{split[j]}, &core.Node{split[j+1]})
				}
			}

			g.String()
			fmt.Println("---------")

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
