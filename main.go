package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func neco(c *gin.Context) {
	bla := c.Param("bla")
	c.JSON(200, gin.H{
		"message": bla,
	})
}

func main() {
	router := gin.Default()

	router.GET("/", neco)

	router.GET("/user/:name", func(c *gin.Context) {
		//name := c.Param("name")
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.GET("/ping/:bla", func(c *gin.Context) {
		x := c.Param("bla")
		c.JSON(200, gin.H{

			"message": x,
		})
	})
	router.Run(":1305")
}
