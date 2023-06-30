package controllers

import "net/http"

func NotFound(template Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		template.Execute(w, r, nil)
	}
}
