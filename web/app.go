package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"hybridAST/modules"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type AppList struct {
	gorm.Model
	AppName   string
	Url       string
	Language  string
	Framework string
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/png":
		saveFile(w, file, handle)
	default:
		jsonResponse(w, http.StatusBadRequest, "El formato de la imágen no es válido")
	}
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}

	err = ioutil.WriteFile("./files/"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	fmt.Println("File uploaded!")
	jsonResponse(w, http.StatusCreated, "Archivo guardado exitosamente")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

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
			//modules.StartAnalyzeZAP()
		}

		if r.Form["arachni"] != nil {
			fmt.Println("Arachni Scan is started")
			modules.StartScanArachni(url)
		}

		fmt.Println("Result:", name, id, framework, url, lang)

		if err := templates.ExecuteTemplate(w, "apps", apptemp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
