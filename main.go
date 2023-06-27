package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	config "runhill.cz/config"
	"runhill.cz/db"
	"runhill.cz/routes"
	"runhill.cz/utils"
)

type Neco struct {
	Email string
}

func Testiky(from int, to int) {

}
func AddOneToEachElement(slice []int) {
	for i := range slice {
		slice[i] += 1920
	}
}

var err error

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(config.ViperConfig)
	//viper.AddConfigPath("/var/www/runhill.cz/config")
	viper.SetConfigType("yml")

	var configuration config.Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	db.Mdb, err = sql.Open(configuration.Database.DBDriver, configuration.Database.DBUser+":"+configuration.Database.DBPassword+"@tcp("+configuration.Database.DBServer+")/"+configuration.Database.DBName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Mdb.Close()

	utils.GoogleClientId = configuration.Google.ClientId
	utils.GoogleClientSecret = configuration.Google.ClientSecret
	utils.GoogleRedirectUrl = configuration.Server.Webname + "/" + configuration.Google.RedirectRoute
	utils.FacebookClientId = configuration.Facebook.ClientId
	utils.FacebookClientSecret = configuration.Facebook.ClientSecret
	utils.FacebookRedirectUrl = configuration.Server.Webname + "/" + configuration.Facebook.RedirectRoute

	utils.SessionName = configuration.Authentication.SessionName
	utils.ServerWebname = configuration.Server.Webname
	utils.SmtpServer = configuration.Email.SmtpServer
	utils.SmtpPort = configuration.Email.SmtpPort
	utils.EmailCharset = configuration.Email.EmailCharset
	utils.EmailFrom = configuration.Email.EmailFrom
	utils.EmailFromName = configuration.Email.EmailFromName
	utils.StaticPath = configuration.Server.StaticPath
	//Testiky(1930, 1990)
	router := routes.NewRouter()

	utils.LoadTemplates(configuration.Server.TemplatePath)
	//utils.LoadTemplates("/var/www/runhill.cz/templates/*.html")
	utils.HttpServer(router, configuration.Server.Port)
	/*
		err = http.ListenAndServe("localhost:1306", routes.New())
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		} */
}
