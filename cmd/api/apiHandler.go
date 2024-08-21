package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joechea-aupp/go-api/cmd/helper"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
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

func (api *Api) getUsers(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if params.ByName("start") == "" {
		params = append(params, httprouter.Param{Key: "start", Value: "0"})
	}

	if params.ByName("limit") == "" {
		params = append(params, httprouter.Param{Key: "limit", Value: "2"})
	}

	start, err := strconv.ParseInt(params.ByName("start"), 10, 64)
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.ParseInt(params.ByName("limit"), 10, 64)
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := api.User.GetUsers(start, limit)
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

func (api *Api) signin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data db.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := api.User.GetUser(data.Username)
	if err != nil {
		helper.ResponseWithError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		helper.ResponseWithError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := helper.GenerateJWT(user.Username, user.Email)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	helper.ResponseWithJSON(w, http.StatusOK, response)
}
