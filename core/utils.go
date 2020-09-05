package core

import (
	"encoding/csv"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

func ImportReport(path string) string {

	jsonFile, err := os.Open(path)
	if err != nil {
		Error.Println("ImportReport:,", err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		Error.Println("ImportReport:", err)
	}

	return string(byteValue)
}

func InitDB() *gorm.DB {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		Error.Println("InitDB:", err)
	}
	//defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Entrypoint{}, &Params{}, &SourceData{})
	db.Model(&Params{}).AddForeignKey("params_id", "entrypoints(id)", "CASCADE", "CASCADE")
	db.Model(&SourceData{}).AddForeignKey("source_name_id", "entrypoints(id)", "CASCADE", "CASCADE")
	return db
}

func GetValue(name gjson.Result, path string) string {
	return gjson.Get(name.String(), path).String()
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
