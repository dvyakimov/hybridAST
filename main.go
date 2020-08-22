package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"hybridAST/core"
	"log"
	"net/http"
	"time"
)

type VulnPageData struct {
	IsNotEmpty bool
	Vulns      []Vuln
}

type Vuln struct {
	Id         uint
	Bug_name   string
	Bug_cwe    string
	Created_at time.Time
	Bug_Path   string
}

var tmpl = template.Must(template.ParseFiles("assets/vulns.html"))

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

	if err := tmpl.ExecuteTemplate(w, "vulns.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "assets/index.html")
}

func main() {

	DBStart()

	r := mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./assets/css/"))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/vulns", Vulnpage)

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
