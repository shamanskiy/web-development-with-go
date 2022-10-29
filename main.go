package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"))
}

func executeTemplate(w http.ResponseWriter, templatePath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Page not found!!!", http.StatusNotFound)
}

func printParamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/txt; charset=utf-8")
	param := chi.URLParam(r, "param")
	fmt.Fprintf(w, "Passed parameter: %s", param)
}

func main() {
	router := chi.NewRouter()

	router.Get("/", homeHandler)
	router.With(middleware.Logger).Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.Get("/print-param/{param}", printParamHandler)
	router.NotFound(http.HandlerFunc(notFoundHandler))

	fmt.Println("Starting a server on :3000")
	http.ListenAndServe(":3000", router)
}
