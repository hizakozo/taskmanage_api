package data

import (
    "github.com/jinzhu/gorm"//ここでパッケージをimport
    "log"
    _ "github.com/go-sql-driver/mysql" //コード内で直接参照するわけではないが、依存関係のあるパッケージには最初にアンダースコア_をつける

)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("mysql", "root:root@tcp(taskmanage-mysql)/taskmanage")
	if(err != nil) {
		log.Fatal(err)
	}
	Db.SingularTable(true)
	Db.LogMode(true)
	return
}

