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

	emailService := models.NewEmailService(cfg.SMTP)

	userMiddleware := middleware.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMiddleware := middleware.CSRF(cfg.CSRF.Key, cfg.CSRF.Secure)

	router.Use(middleware.Logger, csrfMiddleware, userMiddleware.SetUser)

	usersController := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
		ServerAddress:        cfg.Server.Address,
	}
	usersController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS,
		"current-user.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignUp = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	usersController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml"))
	usersController.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml"))
	usersController.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))

	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.CurrentUserHandler)
	})

	router.Get("/signup", usersController.SignUpHandler)
	router.Post("/signup", usersController.CreateUserHandler)
	router.Get("/signin", usersController.SignInHandler)
	router.Post("/signin", usersController.AuthenticateUserHandler)
	router.Post("/signout", usersController.SignOutHandler)
	router.Get("/forgot-pw", usersController.ForgotPasswordHandler)
	router.Post("/forgot-pw", usersController.RequestPasswordResetHandler)
	router.Get("/reset-pw", usersController.NewPasswordHandler)
	router.Post("/reset-pw", usersController.ResetPasswordHandler)

	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	router.Get("/contact", controllers.Static(contactTemplate))

	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	router.Get("/faq", controllers.FAQ(faqTemplate))

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	router.Get("/", controllers.Static(homeTemplate))

	notFoundTemplate := views.Must(views.ParseFS(templates.FS, "not-found.gohtml", "tailwind.gohtml"))
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
