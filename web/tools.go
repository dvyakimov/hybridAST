package web

import (
	"hybridAST/modules"
	"net/http"
)

type toolsdata struct {
	ZapStatus     bool
	ArachniStatus bool
}

func ToolsHandler(w http.ResponseWriter, r *http.Request) {
	var toolsdata toolsdata

	if modules.CheckArachni("http://arachni:7331") != false {
		toolsdata.ArachniStatus = true
	} else {
		toolsdata.ArachniStatus = false
	}

	if modules.CheckZAP("http://zaproxy:8090") != false {
		toolsdata.ZapStatus = true
	} else {
		toolsdata.ZapStatus = false
	}

	if err := templates.ExecuteTemplate(w, "tools", toolsdata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
