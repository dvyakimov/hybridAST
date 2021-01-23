package web

import (
	"fmt"
	"hybridAST/core"
	"hybridAST/modules"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
		return
	}

	//dbhost := os.Getenv("DB_HOST")
	//dbport := os.Getenv("DB_PORT")

	var db = core.InitDB()

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	defer file.Close()

	r.ParseForm()
	tool := r.Form["tool"]

	r.FormValue("severity")
	var SeverityFlag bool

	if r.Form["severityReport"] != nil {
		SeverityFlag = true
	}

	fmt.Println(r.Form["idApp"][0])

	var apptemp core.AppList
	db.Find(&apptemp, "id=?", r.Form["idApp"][0])

	var filename string

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "text/xml":
		filename = saveFileXml(w, file)
		if tool[0] == "OWASP ZAP" {
			//TODO: Сделать проверку на report
			modules.ImportReportZapXml(filename, apptemp.ID, SeverityFlag)
			jsonResponse(w, http.StatusCreated, "Analysing by OWASP ZAP is started")
		} else {
			jsonResponse(w, http.StatusBadRequest, "There is no tool")
			return
		}
	case "application/json": //to edit
		filename = saveFileJson(w, file)
		if tool[0] == "OWASP ZAP" {
			modules.ImportReportZapJson(filename, apptemp.ID, SeverityFlag)
			jsonResponse(w, http.StatusCreated, "Analysing by OWASP ZAP is started")
		} else if tool[0] == "Arachni" {
			modules.ImportReportArachni(filename, apptemp.ID, SeverityFlag)
			jsonResponse(w, http.StatusCreated, "Analysing by Arachni is started")
		} else {
			jsonResponse(w, http.StatusBadRequest, "There is no tool")
			return
		}

	default:
		jsonResponse(w, http.StatusBadRequest, "The report format is invalid")
		return
	}
}

func saveFileJson(w http.ResponseWriter, file multipart.File) string {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return ""
	}

	err = ioutil.WriteFile("/app/files/report.json", data, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return "/app/files/report.json"
}

func saveFileXml(w http.ResponseWriter, file multipart.File) string {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return ""
	}

	err = ioutil.WriteFile("/app/files/report.xml", data, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return "/app/files/report.xml"
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
