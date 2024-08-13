package web

import (
	"net/http"

	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/joechea-aupp/go-api/ui"
	"github.com/julienschmidt/httprouter"
)

type Web struct {
	User *db.UserService
}

func (web *Web) Routes(router *httprouter.Router) {
	app := &Web{
		User: db.NewUserService(),
	}

	// fileserver return http.handler, no need for handlerfunc.
	// http.FS converts the embedded filesystem into an http.FileSystem interface that the http.FilServer can use.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/assets/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/user", app.user)
}
