package main

import (
	"fmt"
	"net/http"

	"github.com/Shamanskiy/lenslocked/controllers"
	"github.com/Shamanskiy/lenslocked/middleware"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/Shamanskiy/lenslocked/templates"
	"github.com/Shamanskiy/lenslocked/views"
	chi_middleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	// time request execution
	router.Use(chi_middleware.Logger, middleware.CSRF)

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	router.Get("/", controllers.Static(homeTemplate))

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	router.Get("/contact", controllers.Static(contactTemplate))

	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	router.Get("/faq", controllers.FAQ(faqTemplate))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersController := controllers.Users{
		UserService: &userService,
	}
	usersController.Templates.SignUp = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	router.Get("/signup", usersController.SignUpHandler)
	router.Post("/signup", usersController.CreateUserHandler)
	router.Get("/signin", usersController.SignInHandler)
	router.Post("/signin", usersController.AuthenticateUserHandler)
	router.Get("/app/me", usersController.CurrentUserHandler)
	router.NotFound(controllers.NotFound)

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe("localhost:3000", router)
}
