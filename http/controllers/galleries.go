package controllers

import (
	"net/http"

	"github.com/Shamanskiy/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		NewGallery Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) NewGalleryFormHandler(w http.ResponseWriter, r *http.Request) {
	data := galleryData(r.FormValue("title"))
	g.Templates.NewGallery.Execute(w, r, data)

}

type GalleryData struct {
	Title string
}

func galleryData(title string) GalleryData {
	return GalleryData{
		Title: title,
	}
}
