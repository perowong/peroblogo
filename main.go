package main

import (
	"fmt"
	"log"

	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/dao"
	"github.com/perowong/peroblogo/routers"
)

func main() {
	db := dao.ConnectMysql()
	defer db.Close()

	r := routers.SetupRouters()
	if err := r.Run(fmt.Sprintf(":%s", conf.C.App.Port)); err != nil {
		log.Fatalln(err.Error())
	}
}
