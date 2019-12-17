package database

import (
	config "billing-gorilla/core"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Db *gorm.DB
}

func ConnectDB() *Database {

	conf := config.New()

	dbUri := conf.DB.DB_USER + ":" + conf.DB.DB_PWD + "@tcp(" + conf.DB.DB_HOST + ":" + strconv.Itoa(conf.DB.DB_PORT) + ")/" + conf.DB.DB_NAME + "?parseTime=true"

	log.Println(dbUri)
	db, err := gorm.Open(conf.DB.DB_DRIVER, dbUri)

	if err != nil {
		panic(err.Error())
	}

	return &Database{
		Db: db,
	}

	// defer db.Close()
}
