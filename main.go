package main

import (
	"fmt"
	"log"

	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/dao"
	"github.com/perowong/peroblogo/routers"
)

// @title Peroblogo Api doc
// @version 1.0
// @contact.name Pero Wong
// @contact.url https://i.overio.space
// @contact.email ynwangpeng@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host https://i.overio.space
// @BasePath /api

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db := dao.ConnectMysql()
	defer db.Close()

	r := routers.SetupRouters()
	if err := r.Run(fmt.Sprintf(":%s", conf.C.App.Port)); err != nil {
		log.Fatalln(err.Error())
	}
}
