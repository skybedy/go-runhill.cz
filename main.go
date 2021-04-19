package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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

	n := time.Now().Year()
	fmt.Println(n)
	/*
		slicex := []int{100: 102}
		for i := range slicex {
			fmt.Println(i)
		}*/

	var buffer [256]int
	slice := buffer[1:100]

	for i := 1920; i < len(slice); i++ {
		slice[i] = int(i)
	}

	//fmt.Println(len(slice))

	//fmt.Println(slice)

	AddOneToEachElement(slice)
	//fmt.Println(slice)

	var str string
	for i := from; i <= to; i++ {

		str += strconv.Itoa(i) + ","

	}

	strx := strings.TrimRight(str, ",")

	arr := strings.Split(strx, ",")
	fmt.Println(arr)

}
func AddOneToEachElement(slice []int) {
	for i := range slice {
		slice[i] += 1920
	}
}

var err error

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	var configuration config.Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	db.Mdb, err = sql.Open(configuration.Database.DBDriver, configuration.Database.DBUser+":"+configuration.Database.DBPassword+"@/"+configuration.Database.DBName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Mdb.Close()

	utils.GoogleClientId = configuration.Google.ClientId
	utils.GoogleClientSecret = configuration.Google.ClientSecret
	utils.GoogleRedirectUrl = "http://localhost:1305/afterlogin"
	utils.SessionName = configuration.Authentication.SessionName
	utils.ServerWebname = configuration.Server.Webname
	utils.SmtpServer = configuration.Email.SmtpServer
	utils.SmtpPort = configuration.Email.SmtpPort
	utils.EmailCharset = configuration.Email.EmailCharset
	utils.EmailFrom = configuration.Email.EmailFrom
	utils.EmailFromName = configuration.Email.EmailFromName
	//Testiky(1930, 1990)
	router := routes.NewRouter()

	utils.LoadTemplates("templates/*.html")
	//utils.LoadTemplates("/var/www/timechip.cz/go-www.timechip.cz/templates/*.html")
	utils.HttpServer(router, configuration.Server.Port)
	/*
		err = http.ListenAndServe("localhost:1306", routes.New())
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		} */
}
