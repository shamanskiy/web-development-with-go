package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Shamanskiy/lenslocked/http/controllers"
	"github.com/Shamanskiy/lenslocked/http/middleware"
	"github.com/Shamanskiy/lenslocked/migrations"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/Shamanskiy/lenslocked/templates"
	"github.com/Shamanskiy/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	pwResetService := &models.PasswordResetService{
		DB: db,
	}

	galleryService := &models.GalleryService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	userMiddleware := middleware.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMiddleware := middleware.CSRF(cfg.CSRF.Key, cfg.CSRF.Secure)

	router.Use(middleware.Logger, csrfMiddleware, userMiddleware.SetUser)

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	notFoundTemplate := views.Must(views.ParseFS(templates.FS, "notFound.gohtml", "tailwind.gohtml"))

	usersController := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
		ServerAddress:        cfg.Server.Address,
	}
	usersController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS,
		"users/currentUser.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignUp = views.Must(views.ParseFS(templates.FS,
		"users/signUp.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"users/signIn.gohtml", "tailwind.gohtml"))
	usersController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"users/forgotPassword.gohtml", "tailwind.gohtml"))
	usersController.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS,
		"users/checkYourEmail.gohtml", "tailwind.gohtml"))
	usersController.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS,
		"users/resetPassword.gohtml", "tailwind.gohtml",
	))

	galleriesController := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesController.Templates.NewGallery = views.Must(views.ParseFS(templates.FS,
		"galleries/newGallery.gohtml", "tailwind.gohtml"))
	galleriesController.Templates.EditGallery = views.Must(views.ParseFS(templates.FS,
		"galleries/editGallery.gohtml", "tailwind.gohtml"))
	galleriesController.Templates.IndexGalleries = views.Must(views.ParseFS(templates.FS,
		"galleries/indexGalleries.gohtml", "tailwind.gohtml"))
	galleriesController.Templates.NotFound = notFoundTemplate

	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.CurrentUserHandler)
	})

	router.Get("/signup", usersController.SignUpFormHandler)
	router.Post("/signup", usersController.SignUpHandler)
	router.Get("/signin", usersController.SignInFormHandler)
	router.Post("/signin", usersController.SignInHandler)
	router.Post("/signout", usersController.SignOutHandler)
	router.Get("/forgot-password", usersController.ForgotPasswordFormHandler)
	router.Post("/forgot-password", usersController.ForgotPasswordHandler)
	router.Get("/reset-password", usersController.NewPasswordFormHandler)
	router.Post("/reset-password", usersController.NewPasswordHandler)

	// this redirects logged-out users to the sign-in page
	router.Route("/galleries", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/new-gallery", galleriesController.NewGalleryFormHandler)
		r.Get("/", galleriesController.IndexGalleriesHandler)
		r.Post("/", galleriesController.NewGalleryHandler)
		r.Get("/{id}/edit", galleriesController.EditGalleryFormHandler)
		r.Post("/{id}/edit", galleriesController.EditGalleryHandler)
	})

	router.Get("/", controllers.Static(homeTemplate))
	router.Get("/faq", controllers.FAQ(faqTemplate))
	router.Get("/contact", controllers.Static(contactTemplate))
	router.NotFound(controllers.NotFound(notFoundTemplate))

	fmt.Printf("Starting a server on %s...\n", cfg.Server.Address)
	http.ListenAndServe(cfg.Server.Address, router)
}

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// TODO: Read the PSQL values from an ENV variable
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: Read the CSRF values from an ENV variable
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// TODO: Read the server values from an ENV variable
	cfg.Server.Address = "localhost:3000"

	return cfg, nil
}
