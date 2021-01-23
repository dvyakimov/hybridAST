package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"hybridAST/core"
	"net/http"
	"os"
	"strconv"
)

func getIdFromRequest(req *http.Request) int {
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	return id
}

var templates = template.Must(template.ParseGlob("assets/*.html"))

func AppHandler(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()

	var apptemp core.AppList
	db.Find(&apptemp, "id=?", id)

	switch r.Method {
	case "GET":
		if err := templates.ExecuteTemplate(w, "apps", apptemp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
