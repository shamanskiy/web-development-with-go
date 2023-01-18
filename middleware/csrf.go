package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRF(next http.Handler) http.Handler {
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)

	return csrfMw(next)
}
