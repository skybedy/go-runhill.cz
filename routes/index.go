package routes

import (
	"net/http"

	"runhill.cz/utils"
)

//var x = sessions.NewCookieStore([]byte(sessionSecret), nil)
//var sessionGlob *sessions.Session

func indexHandler(res http.ResponseWriter, req *http.Request) {

	utils.ExecuteTemplate(res, "index.html", struct {
		Title    string
		MenuData interface{}
		//Etapy []utils.Etapy
	}{
		Title:    "Hlavní strana",
		MenuData: utils.SessionExists(utils.SessionName, req),
		//Etapy: utils.EtapyList(),
	})
}
