package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID        int
	UserID    int
	Title     string
	Published bool
}

type GalleryService struct {
	DB *sql.DB
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
