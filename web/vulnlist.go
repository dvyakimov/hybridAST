package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"net/http"
	"os"
	"time"
)

type Vuln struct {
	Id         uint
	Bug_name   string
	Bug_cwe    string
	Created_at time.Time
	Bug_Path   string
}

type VulnListPage struct {
	IsNotEmpty bool
	Vulns      []Vuln
}

func VulnlistHandler(w http.ResponseWriter, r *http.Request) {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var entrypointTemp []*core.Entrypoint

	var data VulnListPage

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
		data = VulnListPage{Vulns: Vulns, IsNotEmpty: true}
	}

	if err := templates.ExecuteTemplate(w, "vulnlist", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
