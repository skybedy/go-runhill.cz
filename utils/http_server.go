package utils

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	conf "runhill.cz/config"
)

func HttpServer(router *mux.Router, portx string) {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = conf.HTTP_PORT
	}

	server := &http.Server{
		Handler: router,
		//Addr:    "127.0.0.1:" + port,
		Addr: "0.0.0.0:" + portx,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("main: running simple server on port", port)
	if err := server.ListenAndServe(); err != nil {
		//log.Fatal("main: couldn't start simple server: %v\n", err)
		//log.Fatal().Err(err)
	}
}
