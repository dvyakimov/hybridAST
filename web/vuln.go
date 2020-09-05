package web

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"net/http"
	"os"
)

type VulnPageData struct {
	IsNotEmpty  bool
	Params      []core.Params
	Sources     []core.SourceData
	BugName     string
	BugCWE      string
	BugHostPort string
	BugPath     string
}

func VulnHanlder(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	var data VulnPageData

	var entrypoint core.Entrypoint
	db.Find(&entrypoint, "id=?", id)

	var paramsDb []*core.Params
	db.Find(&paramsDb, "params_id=?", id)
	var Params []core.Params
	if paramsDb != nil {
		for i := range paramsDb {
			var ParamTemp core.Params
			ParamTemp.ParamName = paramsDb[i].ParamName
			ParamTemp.ParamValue = paramsDb[i].ParamValue
			Params = append(Params, ParamTemp)
		}
	}

	var sourceDb []*core.SourceData
	db.Find(&sourceDb, "source_name_id=?", id)
	var Sources []core.SourceData
	if sourceDb != nil {
		for i := range sourceDb {
			var SourceTemp core.SourceData
			SourceTemp.Source = sourceDb[i].Source
			SourceTemp.Severity = sourceDb[i].Severity
			SourceTemp.CreatedAt = sourceDb[i].CreatedAt
			SourceTemp.UpdatedAt = sourceDb[i].UpdatedAt
			SourceTemp.SourceName = sourceDb[i].SourceName
			Sources = append(Sources, SourceTemp)
		}
	}
	data = VulnPageData{
		Params:     Params,
		Sources:    Sources,
		BugName:    entrypoint.BugName,
		BugCWE:     entrypoint.BugCWE,
		BugPath:    entrypoint.BugPath,
		IsNotEmpty: true,
	}

	if err := templates.ExecuteTemplate(w, "vuln", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
