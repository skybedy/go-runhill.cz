package routes

import (
	"net/http"

	"runhill.cz/utils"
)

func etapyHandler(res http.ResponseWriter, req *http.Request) {
	//http.Redirect(res, req, "/", http.StatusFound)

	/*
		session, err := utils.SessionStore.Get(req, "index")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(session)
		delete(session.Values, "neco")
		session.Save(res)*/

	utils.ExecuteTemplate(res, "etapa.html", struct {
		Title    string
		MenuData interface{}
	}{
		Title:    "Etapy",
		MenuData: utils.SessionExists(utils.SessionName, req),
	})

}
