package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

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

	err = saveFile(file, fullPathToFile)
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
