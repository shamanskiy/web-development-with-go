package controllers

import (
	"net/http"

	"guthub.com/Shamanskiy/lenslocked/views"
)

func StaticHandler(template views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, nil)
	}
}
