package main

type ImageRepository interface {
	SaveImage(fileName, userID string) error
	GetImages(userID string) ([]string, error)
	DeleteImages(userID string) error
}
