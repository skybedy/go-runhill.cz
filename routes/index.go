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

		//	session, _ := utils.SessionStore.Get(req, utils.SessionName)

		utils.ExecuteTemplate(res, "index.html", struct {
			Title string
			Login interface{}
		}{
			Title: "Hlavní strana",
			Login: utils.SessionExists(utils.SessionName, req),
		})
	}
}
