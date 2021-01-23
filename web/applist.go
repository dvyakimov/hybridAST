package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"net/http"
	"os"
)

type AppPageData struct {
	IsNotEmpty bool
	Apps       []core.AppList
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

	var apptemp []*core.AppList
	db.Find(&apptemp)

	if apptemp != nil {
		var apps []core.AppList
		for i := range apptemp {
			var AppTemp core.AppList
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
	}
}
