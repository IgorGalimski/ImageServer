package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	r := gin.Default()

	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	uploadPath := "uploads"
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	imageRepo := NewLocalImageRepository(uploadPath, db)
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
	defer db.Close()

	return db, nil
}

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "Running!")
}

type ImageHandler struct {
	imageRepo ImageRepository
}

func NewImageHandler(imageRepo ImageRepository) *ImageHandler {
	return &ImageHandler{
		imageRepo: imageRepo,
	}
}

func (h *ImageHandler) uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error getting image from user: %s", err.Error()))
		return
	}

	fileName := generateFileName(file.Filename)
	fullPathToFile := filepath.Join("uploads", fileName)

	err = saveUploadedFile(file, fullPathToFile)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file: %s", err.Error()))
		return
	}

	userId := c.PostForm("userId")
	if userId == "" {
		c.String(http.StatusBadRequest, "Error getting user id")
		return
	}

	err = h.imageRepo.SaveImage(fileName, userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file to database: %s", err.Error()))
		deleteFile(fullPathToFile)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}

func generateFileName(originalName string) string {
	return fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(originalName))
}

func saveUploadedFile(file *multipart.FileHeader, fullPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Println("Error deleting file:", err.Error())
	}
}

func (h *ImageHandler) getImages(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		c.String(http.StatusBadRequest, "Error getting user id")
		return
	}

	filenames, err := h.imageRepo.GetImages(userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error querying database: %s", err.Error()))
		return
	}

	files, err := readFiles(filenames)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading files: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userId,
		"files": files,
	})
}

func readFiles(filenames []string) ([]gin.H, error) {
	var files []gin.H
	for _, filename := range filenames {
		fullPath := filepath.Join("uploads", filename)

		fileBytes, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		files = append(files, gin.H{
			"filename": filename,
			"data":     fileBytes,
		})
	}

	return files, nil
}

func (h *ImageHandler) deleteImages(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		c.String(http.StatusBadRequest, "Error getting user id")
		return
	}

	filenames, err := h.imageRepo.GetImages(userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error querying database: %s", err.Error()))
		return
	}

	err = deleteFiles(filenames)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting files: %s", err.Error()))
		return
	}

	err = h.imageRepo.DeleteImages(userId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting images: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Deleted %d images for user %s", len(filenames), userId))
}

func deleteFiles(filenames []string) error {
	for _, filename := range filenames {
		fullPath := filepath.Join("uploads", filename)
		err := os.Remove(fullPath)
		if err != nil {
			log.Println("Error deleting file:", err.Error())
		}
	}

	return nil
}
