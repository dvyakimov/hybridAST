package web

import (
	"fmt"
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

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	defer file.Close()

	r.ParseForm()
	tool := r.Form["tool"]

	var filename string

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "application/json": //to edit
		filename = saveFile(w, file)
	default:
		jsonResponse(w, http.StatusBadRequest, "The image format is invalid")
		return
	}

	if tool[0] == "OWASP ZAP" {
		modules.ImportReportZap(filename)
		jsonResponse(w, http.StatusCreated, "Analysing by OWASP ZAP is started")
	} else if tool[0] == "Arachni" {
		modules.ImportReportArachni(filename)
		jsonResponse(w, http.StatusCreated, "Analysing by Arachni is started")
	} else {
		jsonResponse(w, http.StatusBadRequest, "There is no tool")
		return
	}
}

func saveFile(w http.ResponseWriter, file multipart.File) string {
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

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
