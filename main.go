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

	r.HandleFunc("/apps", web.AppListHandler)
	r.HandleFunc("/vulns", web.Vulnpage)
	r.HandleFunc("/tools", web.ToolsHandler)
	r.HandleFunc("/apps/{id:[0-9]}", web.AppHandler)

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/apps/css/", http.StripPrefix("/apps/css/", cssHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
