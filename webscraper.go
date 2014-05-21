package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	//For populating struct with multiple url values
)

// IMPORTANT: Make sure that fields are exported, i.e. Capitalized first letters.
type Selection struct {
	Selector string `schema:"s"`
	URL      string `schema:"url"`
}

var templates = template.Must(template.ParseFiles("index.html"))

// InputHandler returns a HTML form for selector string input.
func InputHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// renderTemplate renders the template and handles errors.
// It takes http.Response Writer and the template filename as inputs.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, tmpl+".html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	io.Copy(w, buf)
}

func init() {
	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/result", ResultHandler)
}
