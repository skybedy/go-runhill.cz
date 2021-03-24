package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"runhill.cz/conf"
	"runhill.cz/db"
	"runhill.cz/routes"
	"runhill.cz/utils"
)

func Testiky() {

}

var err error

func main() {
	Testiky()
	db.Mdb, err = sql.Open(conf.DbDriver, conf.DbUser+":"+conf.DbPass+"@/"+conf.DbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Mdb.Close()
	router := routes.NewRouter()
	utils.LoadTemplates("templates/*.html")
	//utils.LoadTemplates("/var/www/timechip.cz/go-www.timechip.cz/templates/*.html")
	utils.HttpServer(router)
}
