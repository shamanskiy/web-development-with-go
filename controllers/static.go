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

func FAQ(template views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Wow, you wrote this yourself?",
			Answer:   "Yes",
		},
		{
			Question: "Can you teach me?",
			Answer:   "Yes",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, questions)
	}
}
