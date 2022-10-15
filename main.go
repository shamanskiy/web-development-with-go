package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Contact Page</h1><p>To get in touch, email me at 
		<a href=\"mailto:megapacha2@gmail.com\">megapacha2@gmail.com</a>.</p>`)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w,
		`<h1>FAQ</h1>
	<ul>
	<li>
	  <b>Wow, you wrote this yourself?</b>
	  Yes!
	</li>
	<li>
	  <b>Can you teach me?</b>
	  Yes
	</li>
  </ul>`)
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
