package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// api endpoint
	app.apiRoutes(router)

	// interface endpoint

	return router
}
