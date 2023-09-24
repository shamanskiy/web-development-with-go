package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shamanskiy/lenslocked/errors"
	"github.com/Shamanskiy/lenslocked/http/context"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		NewGallery     Template
		EditGallery    Template
		IndexGalleries Template
		NotFound       Template
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

func (g Galleries) EditGalleryFormHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		err = errors.Public(err, "Requested gallery with an invalid ID.")
		g.Templates.NotFound.Execute(w, r, struct{}{}, err)
		return
	}
	fmt.Println(id)

	gallery, err := g.GalleryService.FindByID(id)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			err = errors.Public(err, "Requested gallery doesn't exist.")
			g.Templates.NotFound.Execute(w, r, struct{}{}, err)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		err = errors.Public(err, "Requested gallery doesn't exist.")
		g.Templates.NotFound.Execute(w, r, struct{}{}, err)
		return
	}

	g.Templates.EditGallery.Execute(w, r, gallery)
}

func (g Galleries) EditGalleryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		err = errors.Public(err, "Attempted to edit gallery with an invalid ID.")
		g.Templates.NotFound.Execute(w, r, struct{}{}, err)
		return
	}
	fmt.Println(id)

	gallery, err := g.GalleryService.FindByID(id)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			err = errors.Public(err, "Attempted to edit non-existing gallery.")
			g.Templates.NotFound.Execute(w, r, struct{}{}, err)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		err = errors.Public(err, "Attempted to edit non-existing gallery.")
		g.Templates.NotFound.Execute(w, r, struct{}{}, err)
		return
	}

	title := r.FormValue("title")
	gallery.Title = title
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) IndexGalleriesHandler(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Println(user)
	galleries, err := g.GalleryService.FindByUserID(user.ID)
	fmt.Println(err)
	fmt.Println(galleries)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var data struct {
		Galleries []models.Gallery
	}
	data.Galleries = galleries

	g.Templates.IndexGalleries.Execute(w, r, data)
}
