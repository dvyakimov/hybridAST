package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
)

type AppList struct {
	gorm.Model
	AppName   string
	Url       string
	Language  string
	Framework string
}

func getIdFromRequest(req *http.Request) int {
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	return id
}

var templates = template.Must(template.ParseGlob("assets/*.html"))

func AppHandler(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	db, err := gorm.Open("mysql", "root:root@(godb:3306)/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var apptemp AppList
	db.Find(&apptemp, "id=?", id)

	switch r.Method {
	case "GET":
		if err := templates.ExecuteTemplate(w, "apps", apptemp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		name := r.FormValue("appname")
		id := r.FormValue("appid")
		framework := r.FormValue("appframe")
		url := r.FormValue("appurl")
		lang := r.FormValue("applang")

		if r.Form["zaproxy"] != nil {
			fmt.Println("OWASP ZAP Scan is started")

		}

		if r.Form["arachni"] != nil {
			fmt.Println("Arachni Scan is started")

		}

		fmt.Println("Result:", name, id, framework, url, lang)

		if err := templates.ExecuteTemplate(w, "apps", apptemp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
