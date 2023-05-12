package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	r.POST("/upload", uploadImage)
	r.GET("/get", getImages)

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
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error getting image from user: %s", err.Error()))
		return
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))
	fullPathToFile := fmt.Sprintf("uploads/%s", fileName)

	err = c.SaveUploadedFile(file, fullPathToFile)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file: %s", err.Error()))
		return
	}

	f, err := os.Open(fullPathToFile)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error opening file: %s", err.Error()))
		return
	}
	defer f.Close()

	userId := c.PostForm("userId")
	if userId == "" {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error getting user id: %s", err.Error()))
		return
	}

	db := getDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO images (fileName, userId) VALUES (?, ?)")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileName, userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file to database: %s", err.Error()))

		_ = os.Remove(fullPathToFile)

		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}

func getImages(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		c.String(http.StatusBadRequest, "Error getting user id")
		return
	}

	db := getDB()
	defer db.Close()

	rows, err := db.Query("SELECT fileName FROM images WHERE userId = ?", userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error querying database: %s", err.Error()))
		return
	}
	defer rows.Close()

	var files []gin.H
	for rows.Next() {
		var fileName string
		err := rows.Scan(&fileName)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading database results: %s", err.Error()))
			return
		}

		fullPathToFile := fmt.Sprintf("uploads/%s", fileName)

		f, err := os.Open(fullPathToFile)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error opening file: %s", err.Error()))
			return
		}
		defer f.Close()

		fileBytes, err := io.ReadAll(f)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading file: %s", err.Error()))
			return
		}

		files = append(files, gin.H{
			"filename": fileName,
			"data":     fileBytes,
		})
	}

	if err := rows.Err(); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading database results: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userId,
		"files": files,
	})
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/images")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
