package main

import (
	"fmt"
	"net/http"

	"github.com/Shamanskiy/lenslocked/controllers"
	"github.com/Shamanskiy/lenslocked/templates"
	"github.com/Shamanskiy/lenslocked/views"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	router.Get("/", controllers.Static(homeTemplate))

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	router.Get("/contact", controllers.Static(contactTemplate))

	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	router.Get("/faq", controllers.FAQ(faqTemplate))

	var usersController controllers.Users
	usersController.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	router.Get("/signup", usersController.NewHandler)

	router.NotFound(controllers.NotFound)

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe("localhost:3000", router)
}
