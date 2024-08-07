package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// api endpoint
	router.HandlerFunc(http.MethodGet, "/api/healthz", app.healthz)

	// interface endpoint

	return router
}
