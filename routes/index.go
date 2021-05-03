package routes

import (
	"net/http"

	"runhill.cz/utils"
)

//var x = sessions.NewCookieStore([]byte(sessionSecret), nil)
//var sessionGlob *sessions.Session

func indexHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		/*
			session1, err := utils.SessionStore.Get(req, "index")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(session1.Values["lplpl"])
			utils.SessionStore.Destroy(res, "index")*/

	} else {
		/*
			sessionGlob := utils.SessionStore.New("index")
			sessionGlob.Values["neco"] = "blabla"
			sessionGlob.Values["necojin"] = "blplpllablakoko"
			sessionGlob.Save(res)
			fmt.Println(sessionGlob)*/

		utils.ExecuteTemplate(res, "index.html", struct {
			Title string
			Login interface{}
		}{
			Title: "Hlavní strana",
			Login: utils.SessionExists(utils.SessionName, req),
		})
	}
}

func messageHandler(res http.ResponseWriter, req *http.Request) {
	//var message string = ""
	//var alert = "danger"
	//from := req.URL.Query().Get("from")
	//alert = req.URL.Query().Get("alert")

	/*
		switch from {
		case "login":
			message = "Takhle to nepůjde"
		case "loginnoauthorise":
			message = "Účet ještě nebyl funkční, nejprve je třeba dokončit autorizaci, která vám byla zaslána emailem"
		case "authorizationtrue":
			message = "Super, klaplo to a teď se už můžete regulérně přihlásit"
		case "authorizationmisunderstanding":
			message = "Verifikace už byla provedena, pokud nejste, můžete se přihlásit a pokračovat"
		case "accountdelete":
			message = " Váš účet byl smazán, kdykoli se samozřejmě můžete přihlásit znovu."
		case "editpersonsuccess":
			message = "Vaše údaje byly úspěšně změněny."

		}*/

	if message == "" && alert == "" {
		http.Redirect(res, req, "/", http.StatusFound) //tohle se musí pořádně rozpitvat a navrhnout finální řešení, stránka prostě bez parametrů nemá platnost
	}

	utils.ExecuteTemplate(res, "message.html", struct {
		Title     string
		Login     interface{}
		Message   string
		AlertType string
	}{
		Title:     "",
		Login:     utils.SessionExists(utils.SessionName, req),
		Message:   message,
		AlertType: alert,
	})
	message, alert = "", ""
}
