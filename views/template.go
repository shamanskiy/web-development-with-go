package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/Shamanskiy/lenslocked/http/context"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTemplate *template.Template
}

type public interface {
	Public() string
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// cloning template to avoid race condition
	// when handling multiple user requests
	clonedTemplate, err := t.htmlTemplate.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	errorMessage, statusCode := parseError(errs)
	clonedTemplate = clonedTemplate.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errorMessage
		},
	})

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// buffering template execution to avoid sending partial HTML on error
	var buf bytes.Buffer
	err = clonedTemplate.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// placeholder template and csrfField func to parse the template
	htmlTemplate := template.New(patterns[0])
	htmlTemplate = htmlTemplate.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (*models.User, error) {
			return nil, fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return nil
		},
	})

	htmlTemplate, err := htmlTemplate.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTemplate: htmlTemplate}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func parseError(err []error) (errorMessage []string, statusCode int) {
	statusCode = http.StatusOK
	for _, e := range err {
		var publicErr public
		if errors.As(e, &publicErr) {
			statusCode = http.StatusBadRequest
			errorMessage = append(errorMessage, publicErr.Public())
		} else {
			fmt.Println(e)
			statusCode = http.StatusInternalServerError
			errorMessage = append(errorMessage, "Something went wrong.")
		}
	}
	return errorMessage, statusCode
}
