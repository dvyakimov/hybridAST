package main

import (
	"github.com/gorilla/mux"
	"hybridAST/web"
	"log"
	"net/http"
)

func main() {

	DBStart()

	r := mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./assets/css/"))
	imagesHandler := http.FileServer(http.Dir("./assets/images/"))
	jsHandler := http.FileServer(http.Dir("./assets/js/"))

	r.HandleFunc("/apps", web.AppListHandler)
	r.HandleFunc("/upload", web.UploadFile)

	r.HandleFunc("/vulns", web.VulnlistHandler)
	r.HandleFunc("/tools", web.ToolsHandler)
	r.HandleFunc("/apps/{id:[0-9]+}", web.AppHandler)
	r.HandleFunc("/vulns/{id:[0-9]+}", web.VulnHanlder)

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/apps/css/", http.StripPrefix("/apps/css/", cssHandler))
	http.Handle("/vulns/css/", http.StripPrefix("/vulns/css/", cssHandler))

	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))
	http.Handle("/apps/js/", http.StripPrefix("/apps/js/", jsHandler))

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
