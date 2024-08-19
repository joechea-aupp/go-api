package web

import (
	"html/template"
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/middleware"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/joechea-aupp/go-api/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Web struct {
	User          *db.UserService
	templateCache map[string]*template.Template
	templateData  *ui.TemplateData
}

var mid = &middleware.Middleware{}

func (web *Web) Routes(router *httprouter.Router) {
	templateCache, err := ui.NewTemplateCache()
	if err != nil {
		panic(err)
	}

	app := &Web{
		User:          db.NewUserService(),
		templateCache: templateCache,
		templateData:  &ui.TemplateData{},
	}

	webLog := alice.New(mid.LogURL)

	// fileserver return http.handler, no need for handlerfunc.
	// http.FS converts the embedded filesystem into an http.FileSystem interface that the http.FilServer can use.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/assets/*filepath", fileServer)

	router.Handler(http.MethodGet, "/user", webLog.ThenFunc(app.user))
	router.Handler(http.MethodGet, "/count", webLog.ThenFunc(app.count))
	router.Handler(http.MethodPost, "/count/:mode", webLog.ThenFunc(app.postCount))
	router.Handler(http.MethodGet, "/form", webLog.ThenFunc(app.getForm))
	router.Handler(http.MethodPost, "/form", webLog.ThenFunc(app.postForm))
}
