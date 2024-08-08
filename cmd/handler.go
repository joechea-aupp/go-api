package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"` // give it a nickname for json field
	}{
		Status: "ok",
	}
	responseWithJSON(w, http.StatusOK, response)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	// userService := factory.NewUserService()
	params := httprouter.ParamsFromContext(r.Context())
	email := params.ByName("email")

	user, err := app.User.GetUser(email)
	if err != nil {

		errormsg := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprint("error: ", err),
		}
		responseWithJSON(w, http.StatusInternalServerError, errormsg)
		return
	}

	responseWithJSON(w, http.StatusOK, user)
}
