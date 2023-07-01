package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Shamanskiy/lenslocked/errors"
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
	ServerAddress        string
}

func (u Users) SignUpFormHandler(w http.ResponseWriter, r *http.Request) {
	data := emailData(r.FormValue("email"))
	u.Templates.SignUp.Execute(w, r, data)
}

func (u Users) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		u.Templates.SignUp.Execute(w, r, EmailData{}, err)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account.")
		}
		u.Templates.SignUp.Execute(w, r, emailData(email), err)
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

func (u Users) SignInFormHandler(w http.ResponseWriter, r *http.Request) {
	data := emailData(r.FormValue("email"))
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailNotFound) {
			err = errors.Public(err, "No account found associated with this email.")
		} else if errors.Is(err, models.ErrPasswordWrong) {
			err = errors.Public(err, "Provided password is wrong.")
		}
		u.Templates.SignIn.Execute(w, r, emailData(email), err)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, emailData(email), err)
		return
	}
	cookie.Set(w, cookie.CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// This handler expects to sit behind userMiddleware.RequireUser,
// so it doesn't check if the user exists
func (u Users) CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.CurrentUser.Execute(w, r, emailData(user.Email))
}

func (u Users) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := cookie.Read(r, cookie.CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	cookie.Delete(w, cookie.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	data := emailData(r.FormValue("email"))
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	data := emailData(r.FormValue("email"))
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		if errors.Is(err, models.ErrEmailNotFound) {
			err = errors.Public(err, "No account found associated with this email.")
		}
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}
	// TODO: Make the URL here configurable
	err = u.EmailService.ForgotPassword(data.Email,
		"http://"+u.ServerAddress+"/reset-pw?"+vals.Encode())
	if err != nil {
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}

	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) NewPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) NewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		if errors.Is(err, models.ErrInvalidToken) {
			err = errors.Public(err, "Submitted authorization token is invalid.")
		}
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		u.Templates.ResetPassword.Execute(w, r, data, err)
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

type EmailData struct {
	Email string
}

func emailData(email string) EmailData {
	return EmailData{
		Email: email,
	}
}
