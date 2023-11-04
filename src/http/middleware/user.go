package middleware

import (
	"net/http"

	"github.com/Shamanskiy/lenslocked/src/http/context"
	"github.com/Shamanskiy/lenslocked/src/http/cookie"
	"github.com/Shamanskiy/lenslocked/src/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First try to read the cookie. If we run into an error reading it,
		// proceed with the request. The goal of this middleware isn't to limit
		// access. It only sets the user in the context if it can.
		token, err := cookie.Read(r, cookie.CookieSession)
		if err != nil {
			// Cannot lookup the user with no cookie, so proceed without a user being
			// set, then return.
			next.ServeHTTP(w, r)
			return
		}

		// If we have a token, try to lookup the user with that token.
		user, err := umw.SessionService.User(token)
		if err != nil {
			// Invalid or expired token. In either case we can still proceed, we just
			// cannot set a user.
			next.ServeHTTP(w, r)
			return
		}

		// If we get to this point, we have a user that we can store in the context!
		// Get the context
		ctx := r.Context()
		// We need to derive a new context to store values in it. Be certain that
		// we import our own context package, and not the one from the standard
		// library.
		ctx = context.WithUser(ctx, user)
		// Next we need to get a request that uses our new context. This is done
		// in a way similar to how contexts work - we call a WithContext function
		// and it returns us a new request with the context set.
		r = r.WithContext(ctx)
		// Finally we call the handler that our middleware was applied to with the
		// updated request.
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
