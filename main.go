package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"guthub.com/Shamanskiy/lenslocked/controllers"
	"guthub.com/Shamanskiy/lenslocked/templates"
	"guthub.com/Shamanskiy/lenslocked/views"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Page not found!!!", http.StatusNotFound)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	router.Get("/", controllers.StaticHandler(homeTemplate))

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	router.Get("/contact", controllers.StaticHandler(contactTemplate))

	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	router.Get("/faq", controllers.FAQ(faqTemplate))

	router.NotFound(http.HandlerFunc(notFoundHandler))

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe(":3000", router)
}
