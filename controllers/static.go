package controllers

import (
	"net/http"
)

func Static(template Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, r, nil)
	}
}

func FAQ(template Template) http.HandlerFunc {
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
		{
			Question: "What else can I do?",
			Answer:   "Whatever you want",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, r, questions)
	}
}
