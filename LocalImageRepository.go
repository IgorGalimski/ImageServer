package main

import "database/sql"

type LocalImageRepository struct {
	uploadPath string
	db         *sql.DB
}

func NewLocalImageRepository(uploadPath string, db *sql.DB) *LocalImageRepository {
	return &LocalImageRepository{
		uploadPath: uploadPath,
		db:         db,
	}
}
