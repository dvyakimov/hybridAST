package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
)

func AddAppToDb(db *gorm.DB, name string, url string, language string, framework string, contextroot string) {
	app := AppList{
		AppName:     name,
		Language:    language,
		Url:         url,
		Framework:   framework,
		ContextRoot: contextroot,
	}
	db.Create(&app)
}

func AddApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
		return
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// for source code in zip
	//file, handle, err := r.FormFile("file")
	//if err != nil {
	//	fmt.Fprintf(w, "%v", err.Error())
	//	//return
	//}
	//defer file.Close()

	//mimeType := handle.Header.Get("Content-Type")
	//switch mimeType {
	//case "application/zip": //to edit
	//	saveFile(w, file, handle)
	//}

	name := r.FormValue("name")
	url := r.FormValue("url")
	language := r.FormValue("language")
	framework := r.FormValue("framework")
	contextroot := r.FormValue("context-root")

	AddAppToDb(db, name, url, language, framework, contextroot)

	//files, err := Unzip("./filies/" + handle.Filename, "output-folder")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	jsonResponse(w, http.StatusCreated, "Applications is added successfully")
}
