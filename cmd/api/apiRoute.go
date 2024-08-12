package api

import (
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/middleware"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Api struct {
	User *db.UserService
}

var mid = &middleware.Middleware{}

func (api *Api) Routes(router *httprouter.Router) {
	app := &Api{
		User: db.NewUserService(),
	}

	protected := alice.New(mid.VerifyAuth)

	router.HandlerFunc(http.MethodGet, "/healthz", app.healthz)
	router.HandlerFunc(http.MethodGet, "/api/user/:email", app.getUser)
	router.Handler(http.MethodGet, "/api/users", protected.ThenFunc(app.getUsers))
	router.HandlerFunc(http.MethodDelete, "/api/user/:id", app.deleteUser)
	router.HandlerFunc(http.MethodPatch, "/api/user/:id", app.updateUser)

	// user login and register handler
	router.HandlerFunc(http.MethodPost, "/api/register", app.postUser)
	router.HandlerFunc(http.MethodPost, "/api/login", app.signin)
}
