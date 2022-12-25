package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/dao"
	"github.com/perowong/peroblogo/routers"
)

func connectDB() {
	db, err := sqlx.Connect("mysql", conf.C.Mysql["db-peroblog"].Dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(5)
	dao.Setup(db)
}

func startGinApp() {
	if conf.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	routers.Setup(r)

	err := r.Run(fmt.Sprintf(":%s", conf.C.App.Port))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	connectDB()
	startGinApp()
}
