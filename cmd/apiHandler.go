package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joechea-aupp/go-api/internal/db"
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

func (app *application) getUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := app.User.GetUsers()
	if err != nil {
		errormsg := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprint("error: ", err),
		}
		responseWithJSON(w, http.StatusInternalServerError, errormsg)
		return
	}

	responseWithJSON(w, http.StatusOK, users)
}

func (app *application) postUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errormsg := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprint("error: ", err),
		}
		responseWithJSON(w, http.StatusInternalServerError, errormsg)
		return
	}

	defer r.Body.Close()

	var data db.User
	err = json.Unmarshal(body, &data)
	if err != nil {
		errormsg := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprint("error: ", err),
		}
		responseWithJSON(w, http.StatusInternalServerError, errormsg)
		return
	}

	err = app.User.CreateUser(data)
	if err != nil {
		errormsg := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprint("error: ", err),
		}
		responseWithJSON(w, http.StatusInternalServerError, errormsg)
		return
	}

	responseWithJSON(w, http.StatusCreated, "created")
}
