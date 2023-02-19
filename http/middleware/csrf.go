package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRF(csrfKey string, secure bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		csrfMw := csrf.Protect(
			[]byte(csrfKey),
			// TODO: Fix this before deploying
			csrf.Secure(secure),
		)
		return csrfMw(next)
	}
}
