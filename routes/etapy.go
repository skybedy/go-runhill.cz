package routes

import (
	"fmt"
	"net/http"

	"runhill.cz/utils"
)

func etapyHandler(res http.ResponseWriter, req *http.Request) {
	//http.Redirect(res, req, "/", http.StatusFound)

	session, err := utils.SessionStore.Get(req, "index")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(session)
	delete(session.Values, "neco")
	session.Save(res)

}
