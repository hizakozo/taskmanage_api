package infrastructure

import (
	_ "github.com/go-sql-driver/mysql" //コード内で直接参照するわけではないが、依存関係のあるパッケージには最初にアンダースコア_をつける
	"github.com/jinzhu/gorm"           //ここでパッケージをimport
	"log"
	"taskmanage_api/src/constants"
)

var Db *gorm.DB

var url = constants.Params.DbUser + ":" + constants.Params.DbPass + "@tcp(" + constants.Params.DbUrl + ")/taskmanage"

func init() {
	var err error
	Db, err = gorm.Open(
		"mysql",
		url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected !!")
	Db.SingularTable(true)
	Db.LogMode(true)
	return
}
