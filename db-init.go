package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"log"
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

	lines, err := core.ReadCsv("cwe/2000.csv")

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
