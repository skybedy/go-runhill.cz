package main

import (
	"github.com/gin-gonic/gin"
	"runhill.cz/routes"
)

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", routes.Index)
	router.Run(":1305")
}
