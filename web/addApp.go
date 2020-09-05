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

func AddApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	url := r.FormValue("url")
	language := r.FormValue("language")
	framework := r.FormValue("framework")

	AddAppToDb(db, name, url, language, framework)

	fmt.Println(name)

	jsonResponse(w, http.StatusCreated, "Applications is added successfully")
}
