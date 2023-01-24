package dao

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/perowong/peroblogo/conf"
)

var (
	DB *sqlx.DB
)

func ConnectMysql() *sqlx.DB {
	var err error

	DB, err = sqlx.Connect("mysql", conf.C.Mysql["db-peroblog"].Dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(5)

	return DB
}
