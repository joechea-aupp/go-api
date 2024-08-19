package ui

import (
	"html/template"
	"path/filepath"

	"github.com/joechea-aupp/go-api/cmd/middleware"
	"github.com/joechea-aupp/go-api/internal/db"
)

type TemplateData struct {
	User  *db.User
	Users []db.User
	Count int
	URL   map[string]string
	Test  string
}

func NewTemplateCache() (map[string]*template.Template, error) {
	functions := template.FuncMap{
		"currentURL": func() string { return middleware.Feed.Web["path"] },
	}
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
