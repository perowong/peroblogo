package scripts

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/perowong/peroblogo/conf"
)

func GetSqlExecContext() func(sql string) {
	db, err := sqlx.Connect("mysql", conf.C.Mysql["db-peroblog"].Dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return func(sql string) {
		_, err := db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
