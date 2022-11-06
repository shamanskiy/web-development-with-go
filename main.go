package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"guthub.com/Shamanskiy/lenslocked/controllers"
	"guthub.com/Shamanskiy/lenslocked/views"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Page not found!!!", http.StatusNotFound)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	homeTemplate := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	router.Get("/", controllers.StaticHandler(homeTemplate))

	contactTemplate := views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
	router.Get("/contact", controllers.StaticHandler(contactTemplate))

	faqTemplate := views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))
	router.Get("/faq", controllers.StaticHandler(faqTemplate))

	router.NotFound(http.HandlerFunc(notFoundHandler))

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe(":3000", router)
}
