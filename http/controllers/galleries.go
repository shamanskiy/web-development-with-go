package controllers

import (
	"fmt"
	"net/http"
	"sort"
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
		ViewGallery    Template
	}
	GalleryService *models.GalleryService
}

const (
	GALLERY_PUBLIC  = "public"
	GALLERY_PRIVATE = "private"
)

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
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	g.Templates.EditGallery.Execute(w, r, gallery)
}

func (g Galleries) EditGalleryHandler(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	title := r.FormValue("title")
	visibility := r.FormValue("visibility")
	gallery.Title = title
	gallery.Published = visibility == GALLERY_PUBLIC
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
	galleries, err := g.GalleryService.FindByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	sort.Slice(galleries, func(a, b int) bool {
		return galleries[a].ID < galleries[b].ID
	})
	var data struct {
		Galleries []models.Gallery
	}
	data.Galleries = galleries

	g.Templates.IndexGalleries.Execute(w, r, data)
}

func (g Galleries) ViewGalleryHandler(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnPrivateGallery)
	if err != nil {
		return
	}

	var data struct {
		Title  string
		Images []models.Image
	}
	data.Title = gallery.Title
	data.Images, err = g.GalleryService.Images(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	g.Templates.ViewGallery.Execute(w, r, data)
}

func (g Galleries) DeleteGalleryHandler(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(*gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) ImageHandler(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnPrivateGallery)
	if err != nil {
		return
	}
	filename := chi.URLParam(r, "filename")

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var requestedImage models.Image
	imageFound := false
	for _, image := range images {
		if image.Filename == filename {
			requestedImage = image
			imageFound = true
			break
		}
	}

	if !imageFound {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, requestedImage.Path)
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.FindByID(id)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}

	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}

func userMustOwnPrivateGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	if gallery.Published {
		return nil
	}

	user := context.User(r.Context())
	if user == nil || user.ID != gallery.UserID {
		http.Error(w, "You are not authorized to view this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}

	return nil

}
