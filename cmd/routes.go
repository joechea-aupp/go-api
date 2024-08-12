package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// api endpoint
	app.Api.Routes(router)

	// interface endpoint

	standard := alice.New(app.Middleware.LogRequest)

	return standard.Then(router)
}
