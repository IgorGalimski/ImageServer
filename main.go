package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()

	r.POST("/upload", uploadImage)
	r.GET("/", ping)

	err := r.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}

func ping(c *gin.Context) {
	getDB()

	c.String(http.StatusOK, "Running!")
}

func uploadImage(c *gin.Context) {

}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/images")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
