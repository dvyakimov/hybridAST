package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"net/http"
	"time"
)

type Vuln struct {
	Id         uint
	Bug_name   string
	Bug_cwe    string
	Created_at time.Time
	Bug_Path   string
}

type VulnPageData struct {
	IsNotEmpty bool
	Vulns      []Vuln
}

func Vulnpage(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:root@(godb:3306)/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var entrypointTemp []*core.Entrypoint

	var data VulnPageData

	db.Find(&entrypointTemp)
	if entrypointTemp != nil {
		var Vulns []Vuln
		for i := range entrypointTemp {
			var VulnTemp Vuln
			VulnTemp.Id = entrypointTemp[i].ID
			VulnTemp.Bug_name = entrypointTemp[i].BugName
			VulnTemp.Bug_cwe = entrypointTemp[i].BugCWE
			VulnTemp.Created_at = entrypointTemp[i].CreatedAt
			VulnTemp.Bug_Path = entrypointTemp[i].BugPath
			Vulns = append(Vulns, VulnTemp)
		}
		data = VulnPageData{Vulns: Vulns, IsNotEmpty: true}
	}

	if err := templates.ExecuteTemplate(w, "vulnspage", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
