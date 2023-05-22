package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func generateFileName(originalName string) string {
	return fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(originalName))
}

func saveFile(file *multipart.FileHeader, fullPath string) error {
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
