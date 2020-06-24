package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var db *sql.DB

func ConnectDB() {
	if err := ConnectDatabase(); err != nil {
		panic(err)
	}
}

func ConnectDatabase() (err error) {
	if db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"); err != nil {
		return
	}
	err = db.Ping()
	return
}

func InitDB() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	CreateDatabase()
}

func CreateDatabase() {
	_, err := db.Exec(`CREATE DATABASE dbreport`)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(`dbreport successfully created database..`)
	}
	/*
		_,err = db.Exec(`CREATE DATABASE cweDB`)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(`cweDB successfully created database..`)
		}

		_,err = db.Exec(`USE cweDB`)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(`DB selected successfully..`)
		} */
	/*
		stmt, err := db.Prepare(createTable)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(`Table created successfully..`)
		}
	*/
}

var createTable = `
 CREATE TABLE IF NOT EXISTS bugs (
     bug_id          INTEGER PRIMARY KEY AUTO_INCREMENT
     ,bug_name       VARCHAR(32)
     ,bug_cwe        VARCHAR(32)
     ,bug_severity   VARCHAR(32)
     ,bug_url        VARCHAR(32)                    
);
`

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

type CweList struct {
	//ID int64 `gorm:"primary_key"`
	CweID string `gorm:"primary_key"`
	Name  string
}

func DBStart() {
	ConnectDB()
	//InitDB()

	db, err := gorm.Open("mysql", "root:root@/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&CweList{})

	lines, err := ReadCsv("cwe/2000.csv")

	if err != nil {
		fmt.Println(err)
	}

	for _, line := range lines {
		data := CweList{
			CweID: line[0],
			Name:  line[1],
		}
		db.FirstOrCreate(&CweList{}, &data)
	}
	fmt.Println("Done")

}
