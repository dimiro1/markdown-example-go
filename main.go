package main

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/yuin/goldmark"
)

var (
	//go:embed public/**/*
	publicFS embed.FS
	//go:embed templates/**
	templatesFS embed.FS
)

func main() {
	// Parse Templates
	templates := template.Must(
		template.New("").
			ParseFS(templatesFS, "templates/*.gohtml"),
	)

	// Serve static files
	http.Handle("/public/", http.FileServer(http.FS(publicFS)))

	// Setup handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = templates.ExecuteTemplate(w, "index.gohtml", nil)
	})

	http.HandleFunc("/markdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "not able to parse the form", http.StatusBadRequest)
			return
		}

		var (
			message = r.FormValue("message")
			output  = &bytes.Buffer{}
		)

		if err := goldmark.Convert([]byte(message), output); err != nil {
			http.Error(w, "not able to process the markdown", http.StatusInternalServerError)
			return
		}

		if err := templates.ExecuteTemplate(w, "markdown.gohtml", map[string]interface{}{
			"Message": template.HTML(output.String()),
		}); err != nil {
			http.Error(w, "not able to process the markdown", http.StatusInternalServerError)
			return
		}
	})

	// Start the server
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}
