package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"

	"github.com/gorilla/schema"
)

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	var s = new(Selection) //Returns a pointer to a new Selection type
	var decoder = schema.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = decoder.Decode(s, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.Selector == "" || s.URL == "" {
		http.Redirect(w, r, "/", http.StatusFound) //If the query or the format is empty, redirect to the home page.
		return
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	result := matchSelector(s, client)

	fmt.Fprintf(w, "%s : %s\n-----\n%s", s.URL, s.Selector, result)
}

func matchSelector(s *Selection, client *http.Client) string {
	link, err := url.Parse(s.URL)
	if err != nil {
		log.Fatal("Incorrect url")
		return ""
	}
	r, err := client.Get(link.String())
	if err != nil {
		log.Fatal(err)

	}
	doc, err := html.Parse(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	sel, err := cascadia.Compile(s.Selector)
	if err != nil {
		log.Fatal(err)
	}
	matches := sel.MatchAll(doc)
	var result string
	for _, m := range matches {
		result += nodeString(m)
		result += "\n"
	}
	return result
}

func nodeString(n *html.Node) string {
	switch n.Type {
	case html.TextNode:
		return n.Data
	case html.ElementNode:
		return html.Token{
			Type: html.StartTagToken,
			Data: n.Data,
			Attr: n.Attr,
		}.String()
	}
	return ""
}
