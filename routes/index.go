package routes

import (
	"net/http"

	"runhill.cz/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "index.html", struct {
		Title string
	}{
		Title: "Hlavní strana",
	})
}
