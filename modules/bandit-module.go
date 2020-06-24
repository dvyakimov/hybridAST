package modules

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"hybridAST/entrypoint"
)

func StartAnalyzeBandit() {
	db, err := gorm.Open("mysql", "root:root@/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&entrypoint.Entrypoint{}, &entrypoint.Params{}, &entrypoint.SourceData{})
	db.Model(&entrypoint.Params{}).AddForeignKey("params_id", "entrypoints(id)", "CASCADE", "CASCADE")
	db.Model(&entrypoint.SourceData{}).AddForeignKey("source_name_id", "entrypoints(id)", "CASCADE", "CASCADE")

	result := gjson.Get(entrypoint.ImportReport("examples-report/bandit-django.json"), "results")

	for _, name := range result.Array() {
		fmt.Println(gjson.Get(name.String(), "line_number").String())

	}
}
