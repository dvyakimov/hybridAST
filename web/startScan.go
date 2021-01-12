package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/modules"
	"net/http"
	"os"
)

func checkConnect(url string) bool {

	resp, err := http.Get(url)
	if err != nil {
	}
	//defer resp.Body.Close()

	if resp != nil {
		return true
	} else {
		return false
	}
}

func StartScan(w http.ResponseWriter, r *http.Request) {
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

	r.FormValue("zaproxy")
	r.FormValue("arachni")

	var apptemp AppList
	db.Find(&apptemp, "id=?", r.Form["idApp"][0])

	// Сделать проверку, что, если localhost или 127.0.0.1, то менять на host.docker.internal

	if r.Form["zaproxy"] != nil {
		if checkConnect(apptemp.Url) != false {
			fmt.Println("OWASP ZAP Scan is started", r.Form["idApp"][0])
			modules.StartScanZap(apptemp.Url)
		} else {
			fmt.Println("No connection with: ", apptemp.Url)
		}
	}
	if r.Form["arachni"] != nil {
		if checkConnect(apptemp.Url) != false {
			fmt.Println("Arachni Scan is started", r.Form["idApp"][0])
			modules.StartScanArachni(apptemp.Url)
		} else {
			fmt.Println("No connection with: ", apptemp.Url)
		}
	}
}
