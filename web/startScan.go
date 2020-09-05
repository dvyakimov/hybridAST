package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/modules"
	"net/http"
	"os"
)

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
	r.FormValue("arachi")

	var apptemp AppList
	db.Find(&apptemp, "id=?", r.Form["idApp"][0])

	if r.Form["zaproxy"] != nil {
		modules.StartScanZap(apptemp.Url)
		fmt.Println("OWASP ZAP Scan is started", r.Form["idApp"][0])
	}

	if r.Form["arachni"] != nil {
		modules.StartScanArachni(apptemp.Url)
		fmt.Println("Arachni Scan is started", r.Form["idApp"][0])
	}
}
