package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) apiRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/healthz", app.healthz)
	router.HandlerFunc(http.MethodGet, "/api/user/:email", app.getUser)
	router.HandlerFunc(http.MethodGet, "/api/users", app.getUsers)
	router.HandlerFunc(http.MethodPost, "/api/user", app.postUser)
}
