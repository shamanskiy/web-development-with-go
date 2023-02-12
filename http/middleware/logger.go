package middleware

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

var Logger func(next http.Handler) http.Handler = middleware.Logger
