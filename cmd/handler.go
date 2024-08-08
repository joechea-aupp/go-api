package main

import (
	"fmt"
	"net/http"

	"github.com/joechea-aupp/go-api/internal/factory"
	"github.com/julienschmidt/httprouter"
)

var User factory.User

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"` // give it a nickname for json field
	}{
		Status: "ok",
	}
	responseWithJSON(w, http.StatusOK, response)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	email := params.ByName("email")

	user, err := User.GetUser(email)
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
