package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
)

type Entrypoint struct {
	gorm.Model
	BugName     string
	BugCWE      string
	BugHostPort string
	BugPath     string
	Is_correl   bool
	Params      []Params     `gorm:"foreignkey:ParamsID"`
	SourceData  []SourceData `gorm:"foreignkey:SourceNameID"`
}

type Params struct {
	gorm.Model
	ParamName  string
	ParamValue string
	ParamsID   uint
}

type SourceData struct {
	gorm.Model
	SourceNameID uint
	Source       string
	Severity     string
	SourceName   string
	CodeLine     string
	CodeFile     string
	//BugURL string
}

type CweList struct {
	CweID string `gorm:"primary_key"`
	Name  string
}

func FindNameByCWE(db *gorm.DB, NameFromJson string, BugCWE string) string {
	var CweResult []*CweList
	db.Find(&CweResult, "cwe_id=?", BugCWE)
	if len(CweResult) > 0 {
		return CweResult[0].Name
	} else {
		return NameFromJson
	}
}

func UpdateEntry(entry Entrypoint, source SourceData, db *gorm.DB, m url.Values) {

	var entrypointTemp []*Entrypoint
	db.Find(&entrypointTemp, "bug_cwe=? and bug_path=?", entry.BugCWE, entry.BugPath)
	//m := UrlExtractParametr(BugUrl)
	var checkParam = false
	var checkSource = false
	var checkName = false
	var ParamSaveId uint
	if len(entrypointTemp) > 0 {
		for entrypointTempList := 0; entrypointTempList < len(entrypointTemp); entrypointTempList++ {
			var ParamTemp []*Params
			db.Find(&ParamTemp, "params_id=?", entrypointTemp[entrypointTempList].ID)
			if len(ParamTemp) > 0 {
				for ParamTempList := 0; ParamTempList < len(ParamTemp); ParamTempList++ {
					if m != nil {
						for k, _ := range m {
							if ParamTemp[ParamTempList].ParamName == k {
								checkParam = true
								ParamSaveId = ParamTemp[ParamTempList].ParamsID
								break
							} else {
								checkParam = false
							}
						}
					} else {
						if ParamTemp[ParamTempList].ParamName == "" {
							ParamSaveId = ParamTemp[ParamTempList].ParamsID
							checkParam = true
							break
						} else {
							checkParam = false
						}
					}
				}
			}
			if checkParam == true {
				break
			}
		}
		if checkParam == false {
			db.Create(&entry)
			CreateParam(m, db, entry)
			source.SourceNameID = entry.ID
			db.Create(&source)
		}
	} else {
		db.Create(&entry)
		CreateParam(m, db, entry)
		source.SourceNameID = entry.ID
		db.Create(&source)
	}
	var SourceTemp []*SourceData
	db.Find(&SourceTemp, "source_name_id=?", ParamSaveId)
	if len(SourceTemp) > 0 {
		checkSource = CheckSourceResult(SourceTemp, source.Source)
		checkName = CheckSourceNameResult(SourceTemp, source.SourceName)
	}
	if checkParam == true && checkSource == true && checkName == false {
		source.SourceNameID = ParamSaveId
		db.Create(&source)
	}
	if checkParam == true && checkSource == false {
		fmt.Println("Correalted Item")
		source.SourceNameID = ParamSaveId
		db.Create(&source)
		//db.Model(&entry).Update("is_correl","true")
	}
}

func CheckSourceNameResult(SourceTemp []*SourceData, sourceName string) bool {
	var checkName = false
	for SourceTempList := 0; SourceTempList < len(SourceTemp); SourceTempList++ {
		if sourceName == SourceTemp[SourceTempList].SourceName {
			checkName = true
			break
		} else {
			checkName = false
		}
	}
	return checkName
}

func CheckSourceResult(SourceTemp []*SourceData, source string) bool {
	var checkSource = false
	for SourceTempList := 0; SourceTempList < len(SourceTemp); SourceTempList++ {
		if SourceTemp[SourceTempList].Source == source {
			checkSource = true
			break
		} else {
			checkSource = false
		}
	}
	return checkSource
}

func CreateParam(m url.Values, db *gorm.DB, entry Entrypoint) {
	if m != nil {
		for k, _ := range m {
			param := Params{
				ParamName: k,
				ParamsID:  entry.ID,
			}
			db.Create(&param)
		}
	} else {
		param := Params{
			ParamName: "",
			ParamsID:  entry.ID,
		}
		db.Create(&param)
	}
}

func UrlExtractPath(urlFull string) string {
	u, err := url.Parse(urlFull)
	if err != nil {
		log.Println(err)
	}
	return u.Path
}

func UrlExtractHostPort(urlFull string) string {
	u, err := url.Parse(urlFull)
	if err != nil {
		log.Println(err)
	}
	return u.Host
}
