package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/login", Login).Methods("GET")
	router.HandleFunc("/filetesty", Filetesty).Methods("POST")
	//staticFileDirectory := http.Dir("/var/www/timechip.cz/go-www.timechip.cz/static")
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	return router
}
