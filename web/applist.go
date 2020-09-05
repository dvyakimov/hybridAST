package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
)

func AddAppToDb(db *gorm.DB, name string, url string, language string, framework string) {
	app := AppList{
		AppName:   name,
		Language:  language,
		Url:       url,
		Framework: framework,
	}
	db.Create(&app)
}

type AppPageData struct {
	IsNotEmpty bool
	Apps       []AppList
}

func AppListHandler(w http.ResponseWriter, r *http.Request) {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
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
		if err := templates.ExecuteTemplate(w, "applist", data); err != nil {
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
		if err := templates.ExecuteTemplate(w, "applist", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
