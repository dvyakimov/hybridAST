package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"hybridAST/core"
	"log"
	"net/http"
	"strconv"
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

type AppPageData struct {
	IsNotEmpty bool
	Apps       []AppList
}

var templates = template.Must(template.ParseGlob("assets/*.html"))

func Vulnpage(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:root@(localhost:3306)/dbreport?charset=utf8&parseTime=True&loc=Local")
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

func AddAppToDb(db *gorm.DB, name string, url string, language string, framework string) {
	app := AppList{
		AppName:   name,
		Language:  language,
		Url:       url,
		Framework: framework,
	}
	db.Create(&app)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:root@(localhost:3306)/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var data AppPageData

	var apptemp []*AppList
	db.Find(&apptemp)

	if apptemp != nil {
		var apps []AppList
		for i := range apptemp {
			var AppTemp AppList
			AppTemp.ID = apptemp[i].ID
			AppTemp.Language = apptemp[i].Language
			AppTemp.Framework = apptemp[i].Framework
			AppTemp.AppName = apptemp[i].AppName
			AppTemp.Url = apptemp[i].Url
			apps = append(apps, AppTemp)
		}
		data = AppPageData{Apps: apps, IsNotEmpty: true}
	}

	switch r.Method {
	case "GET":
		if err := templates.ExecuteTemplate(w, "homepage", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		name := r.FormValue("name")
		url := r.FormValue("url")
		language := r.FormValue("language")
		framework := r.FormValue("framework")

		AddAppToDb(db, name, url, language, framework)

		fmt.Println(name)
		if err := templates.ExecuteTemplate(w, "homepage", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getIdFromRequest(req *http.Request) int {
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	return id
}

func AppHandler(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	db, err := gorm.Open("mysql", "root:root@(localhost:3306)/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var apptemp AppList
	db.Find(&apptemp, "id=?", id)

	if err := templates.ExecuteTemplate(w, "apps", apptemp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {

	DBStart()

	r := mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./assets/css/"))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/vulns", Vulnpage)
	r.HandleFunc("/apps/{id:[0-9]}", AppHandler)

	http.Handle("*/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/apps/css/", http.StripPrefix("/apps/css/", cssHandler))

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
