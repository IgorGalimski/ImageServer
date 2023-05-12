package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", uploadImage)
	r.GET("/", ping)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Running!")
}

func uploadImage(c *gin.Context) {

}
