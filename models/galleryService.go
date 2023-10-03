package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var supportedExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

type Gallery struct {
	ID        int
	UserID    int
	Title     string
	Published bool
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryService where to store and locate
	// images. If not set, the GalleryService will default to using the "images"
	// directory.
	ImagesDir string
}

type Image struct {
	Path     string
	Filename string
}

func (gs *GalleryService) Create(userId int, title string) (*Gallery, error) {
	gallery := Gallery{
		UserID: userId,
		Title:  title,
	}

	row := gs.DB.QueryRow(`
	  INSERT INTO galleries (user_id, title, published)
	  VALUES ($1, $2, false) RETURNING id`,
		gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID)

	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) FindByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}

	row := gs.DB.QueryRow(`
	  SELECT user_id, title, published
	  FROM galleries WHERE id=$1`, id)
	err := row.Scan(&gallery.UserID, &gallery.Title, &gallery.Published)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, fmt.Errorf("find gallery by id: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) FindByUserID(userId int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
	  SELECT id, title, published
	  FROM galleries WHERE user_id=$1`, userId)
	if err != nil {
		return nil, fmt.Errorf("find galleries by user_id: %w", err)
	}

	galleries := []Gallery{}
	for rows.Next() {
		gallery := Gallery{
			UserID: userId,
		}
		err := rows.Scan(&gallery.ID, &gallery.Title, &gallery.Published)
		if err != nil {
			return nil, fmt.Errorf("find galleries by user_id: %w", err)
		}
		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("find galleries by user_id: %w", rows.Err())
	}

	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
	  UPDATE galleries 
	  SET title=$1, published=$2
		WHERE id=$3`, gallery.Title, gallery.Published, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(gallery Gallery) error {
	_, err := gs.DB.Exec(`
	  DELETE FROM galleries 
		WHERE id=$1`, gallery.ID)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}

func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	var images []Image
	for _, filePath := range allFiles {
		if hasExtension(filePath, supportedExtensions) {
			images = append(images, Image{
				Path:     filePath,
				Filename: filepath.Base(filePath)})
		}
	}

	return images, nil
}

func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(service.galleryDir(galleryID), filename)

	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrImageNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}

	return Image{
		Filename: filename,
		Path:     imagePath,
	}, nil
}

func (service GalleryService) galleryDir(galleryID int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", galleryID))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
