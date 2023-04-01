package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go-db-tool/tool"
)

func main() {
	/*
	   WithDsn() -- Specify dsn (required)
	   WithTable() -- Specify a table name (required)
	   WithRealTableName()  -- Obtain the real table name (not necessary)
	   WithTagKey()  -- Tag tag key, default to gorm tag (not necessary)
	   WithSavePath()  -- Write file path, do not write and print to the console (not necessary)
	*/

	// example
	options := []tool.Options{
		tool.WithDsn("username:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=true&loc=Local"),
		tool.WithTable("table_name"),
		tool.WithRealTableName("real_table_name"),
		tool.WithTagKey("gorm"),
		tool.WithSavePath("./table_name.go"),
	}

	err := tool.NewMysqlTableGoStruct(options...).Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}
