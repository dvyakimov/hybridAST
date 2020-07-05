package core

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

func ImportReport(RerortURL string) string {
	jsonFile, err := os.Open(RerortURL)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return string(byteValue)
}

func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	// Migrate the schema
	db.AutoMigrate(&Entrypoint{}, &Params{}, &SourceData{})
	db.Model(&Params{}).AddForeignKey("params_id", "entrypoints(id)", "CASCADE", "CASCADE")
	db.Model(&SourceData{}).AddForeignKey("source_name_id", "entrypoints(id)", "CASCADE", "CASCADE")
	return db
}

func GetValue(name gjson.Result, path string) string {
	return gjson.Get(name.String(), path).String()
}

func ExtractVaule(name gjson.Result, path string) gjson.Result {
	return gjson.Get(name.String(), path)
}

func ReadCsv(filename string) ([][]string, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return lines, nil
}
