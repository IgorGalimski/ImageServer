package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	uploadPath := "uploads"
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	imageRepo := NewDBImageRepository(uploadPath, db)
	imageHandler := NewImageHandler(imageRepo)

	r.POST("/upload", imageHandler.uploadImage)
	r.GET("/get", imageHandler.getImages)
	r.DELETE("/delete", imageHandler.deleteImages)

	r.GET("/", pingHandler)

	err = r.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/images")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "Running!")
}
