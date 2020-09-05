package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"hybridAST/core"
	"log"
	"os"
)

type CweList struct {
	//ID int64 `gorm:"primary_key"`
	CweID string `gorm:"primary_key"`
	Name  string
}

type AppList struct {
	gorm.Model
	AppName   string
	Url       string
	Language  string
	Framework string
}

var db *sql.DB

func ConnectDB() {
	if err := ConnectDatabase(); err != nil {
		panic(err)
	}
}

func ConnectDatabase() (err error) {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	if db, err = sql.Open("mysql", "root:root@tcp("+dbhost+":"+dbport+")/"); err != nil {
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
	//fmt.Println("Hello 1")
	_, err := db.Exec(`CREATE DATABASE dbreport`)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(`dbreport successfully created database..`)
	}
	//fmt.Println("Hello 2")
}

func DBStart() {
	ConnectDB()
	InitDB()

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	dbGorm, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err)
	}
	defer dbGorm.Close()

	dbGorm.AutoMigrate(&CweList{}, &AppList{})
	//dbGorm.AutoMigrate(&AppList{})

	lines, err := core.ReadCsv("cwe/2000.csv")

	if err != nil {
		fmt.Println(err)
	}

	for _, line := range lines {
		data := CweList{
			CweID: line[0],
			Name:  line[1],
		}
		dbGorm.FirstOrCreate(&CweList{}, &data)
	}
	fmt.Println("Done")

}
