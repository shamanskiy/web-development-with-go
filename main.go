package main

import (
	"fmt"
	"net"
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
	galleriesController.Templates.ViewGallery = views.Must(views.ParseFS(templates.FS,
		"galleries/viewGallery.gohtml", "tailwind.gohtml"))

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
		r.Get("/{id}", galleriesController.ViewGalleryHandler)
		r.Get("/{id}/images/{filename}", galleriesController.ImageHandler)
		r.Group(func(r chi.Router) {
			r.Use(userMiddleware.RequireUser)
			r.Get("/new-gallery", galleriesController.NewGalleryFormHandler)
			r.Get("/", galleriesController.IndexGalleriesHandler)
			r.Post("/", galleriesController.NewGalleryHandler)
			r.Get("/{id}/edit", galleriesController.EditGalleryFormHandler)
			r.Post("/{id}/edit", galleriesController.EditGalleryHandler)
			r.Post("/{id}/delete", galleriesController.DeleteGalleryHandler)
			r.Post("/{id}/images/{filename}/delete", galleriesController.DeleteImageHandler)
			r.Post("/{id}/images", galleriesController.UploadImageHandler)
		})
	})

	router.Get("/", controllers.Static(homeTemplate))
	router.Get("/faq", controllers.FAQ(faqTemplate))
	router.Get("/contact", controllers.Static(contactTemplate))
	router.NotFound(controllers.NotFound(notFoundTemplate))

	assetsHandler := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	fmt.Printf("Listening on http://localhost%s\n", cfg.Server.Address)
	fmt.Printf("Listening on http://%s%s\n", localIpAddress(), cfg.Server.Address)
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

	// read .env file and set env variables
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.User == "" || cfg.PSQL.Database == "" {
		return cfg, fmt.Errorf("no postgres config provided")
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func localIpAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("failed to find local ipv4 address")
}
