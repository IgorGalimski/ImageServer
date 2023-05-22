package main

import (
	"database/sql"
	"path/filepath"
)

type DBImageRepository struct {
	uploadPath string
	db         *sql.DB
}

func NewDBImageRepository(uploadPath string, db *sql.DB) *DBImageRepository {
	return &DBImageRepository{
		uploadPath: uploadPath,
		db:         db,
	}
}

func (r *DBImageRepository) SaveImage(fileName, userID string) error {
	fullPath := filepath.Join(r.uploadPath, fileName)

	stmt, err := r.db.Prepare("INSERT INTO images (fileName, userId) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileName, userID)
	if err != nil {
		deleteFile(fullPath)
		return err
	}

	return nil
}

func (r *DBImageRepository) GetImages(userID string) ([]string, error) {
	var filenames []string

	rows, err := r.db.Query("SELECT fileName FROM images WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		err := rows.Scan(&filename)
		if err != nil {
			return nil, err
		}

		filenames = append(filenames, filename)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return filenames, nil
}

func (r *DBImageRepository) DeleteImages(userID string) error {
	stmt, err := r.db.Prepare("DELETE FROM images WHERE userId = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}
