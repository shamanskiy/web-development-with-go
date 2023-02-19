package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Shamanskiy/lenslocked/http/context"
	"github.com/Shamanskiy/lenslocked/http/cookie"
	"github.com/Shamanskiy/lenslocked/models"
)

type Users struct {
	Templates struct {
		CurrentUser    Template
		SignUp         Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignUp.Execute(w, r, data)
}

func (u Users) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Can't parse the submitted form", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie.Set(w, cookie.CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Can't parse the submitted form", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	cookie.Set(w, cookie.CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// This handler expects to sit behind userMiddleware.RequireUser,
// so it doesn't check if the user exists
func (u Users) CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	var data struct {
		Email string
	}
	data.Email = user.Email
	u.Templates.CurrentUser.Execute(w, r, data)
}

func (u Users) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := cookie.Read(r, cookie.CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	cookie.Delete(w, cookie.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance,
		// if a user doesn't exist with the email address.
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}
	// TODO: Make the URL here configurable
	err = u.EmailService.ForgotPassword(data.Email, "https://www.lenslocked.com/reset-pw?"+vals.Encode())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) NewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between server errors and invalid token errors.
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Sign the user in now that they have reset their password.
	// Any errors from this point onward should redirect to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	cookie.Set(w, cookie.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
