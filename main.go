package main

import (
	"fmt"
	"net/http"

	"github.com/Shamanskiy/lenslocked/http/controllers"
	"github.com/Shamanskiy/lenslocked/http/middleware"
	"github.com/Shamanskiy/lenslocked/migrations"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/Shamanskiy/lenslocked/templates"
	"github.com/Shamanskiy/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	userMiddleware := middleware.UserMiddleware{
		SessionService: &sessionService,
	}

	router.Use(middleware.Logger, middleware.CSRF, userMiddleware.SetUser)

	usersController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS,
		"currentuser.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignUp = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))

	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.CurrentUserHandler)
	})

	router.Get("/signup", usersController.SignUpHandler)
	router.Post("/signup", usersController.CreateUserHandler)
	router.Get("/signin", usersController.SignInHandler)
	router.Post("/signin", usersController.AuthenticateUserHandler)
	router.Post("/signout", usersController.SignOutHandler)

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	router.Get("/contact", controllers.Static(contactTemplate))

	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	router.Get("/faq", controllers.FAQ(faqTemplate))

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	router.Get("/", controllers.Static(homeTemplate))

	router.NotFound(controllers.NotFound)

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe("localhost:3000", router)
}
