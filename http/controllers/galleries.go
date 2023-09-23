package controllers

import (
	"fmt"
	"net/http"

	"github.com/Shamanskiy/lenslocked/http/context"
	"github.com/Shamanskiy/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		NewGallery Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) NewGalleryFormHandler(w http.ResponseWriter, r *http.Request) {
	gallery := models.Gallery{Title: r.FormValue("title")}
	g.Templates.NewGallery.Execute(w, r, gallery)
}

func (g Galleries) NewGalleryHandler(w http.ResponseWriter, r *http.Request) {
	userID := context.User(r.Context()).ID
	title := r.FormValue("title")

	gallery, err := g.GalleryService.Create(userID, title)
	if err != nil {
		g.Templates.NewGallery.Execute(w, r, gallery, err)
		return
	}

	editGalleryPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editGalleryPath, http.StatusFound)
}
