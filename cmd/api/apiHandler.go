package api

import (
	"encoding/json"
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/helper"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/julienschmidt/httprouter"
)

func (api *Api) healthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"` // give it a nickname for json field
	}{
		Status: "ok",
	}
	helper.ResponseWithJSON(w, http.StatusOK, response)
}

func (api *Api) getUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	email := params.ByName("email")

	user, err := api.User.GetUser(email)
	if err != nil {
		helper.ResponseWithError(w, http.StatusNotFound, "user not found")
		return
	}

	helper.ResponseWithJSON(w, http.StatusOK, user)
}

func (api *Api) getUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := api.User.GetUsers()
	if err != nil {
		helper.ResponseWithError(w, http.StatusNotFound, "user not found")
		return
	}

	helper.ResponseWithJSON(w, http.StatusOK, users)
}

func (api *Api) postUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data db.User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = api.User.CreateUser(data)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	helper.ResponseWithJSON(w, http.StatusCreated, "created")
}

func (api *Api) deleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := api.User.DelUser(id)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	helper.ResponseWithJSON(w, http.StatusOK, "deleted")
}

func (api *Api) updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var data db.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := api.User.UpdateUser(id, data); err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	helper.ResponseWithJSON(w, http.StatusOK, "updated")
}
