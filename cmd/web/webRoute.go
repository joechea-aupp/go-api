package web

import (
	"html/template"
	"net/http"

	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/joechea-aupp/go-api/ui"
	"github.com/julienschmidt/httprouter"
)

type Web struct {
	User          *db.UserService
	Form          *Form
	templateCache map[string]*template.Template
}

type Form struct {
	FirstName string
	LastName  string
	Count     int
}

func (web *Web) Routes(router *httprouter.Router) {
	templateCache, err := newTemplateCache()
	if err != nil {
		panic(err)
	}

	app := &Web{
		User:          db.NewUserService(),
		Form:          &Form{},
		templateCache: templateCache,
	}

	// fileserver return http.handler, no need for handlerfunc.
	// http.FS converts the embedded filesystem into an http.FileSystem interface that the http.FilServer can use.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/assets/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/user", app.user)
	router.HandlerFunc(http.MethodGet, "/count", app.count)
	router.HandlerFunc(http.MethodPost, "/count", app.postCount)
	router.HandlerFunc(http.MethodGet, "/form", app.getForm)
	router.HandlerFunc(http.MethodPost, "/form", app.postForm)
}
