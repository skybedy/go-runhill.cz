package routes

import (
	"fmt"
	"net/http"

	"runhill.cz/utils"
)

//var x = sessions.NewCookieStore([]byte(sessionSecret), nil)

func indexHandler(w http.ResponseWriter, req *http.Request) {

	/*
		type Cookie struct {
			Name       string
			Value      string
			Path       string
			Domain     string
			Expires    time.Time
			RawExpires string

			// MaxAge=0 means no 'Max-Age' attribute specified.
			// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
			// MaxAge>0 means Max-Age attribute present and given in seconds
			MaxAge   int
			Secure   bool
			HttpOnly bool
			Raw      string
			Unparsed []string // Raw text of unparsed attribute-value pairs
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "test", Value: "testovič", Expires: expiration}
		http.SetCookie(w, &cookie)*/

	utils.ExecuteTemplate(w, "index.html", struct {
		Title string
		Login interface{}
	}{
		Title: "Hlavní strana",
		Login: utils.SessionExists(utils.SessionName, req),
	})
}

func messageHandler(res http.ResponseWriter, req *http.Request) {
	var message string = ""
	var alert = "danger"
	from := req.URL.Query().Get("from")
	alert = req.URL.Query().Get("alert")
	fmt.Println(len(alert))

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
}
