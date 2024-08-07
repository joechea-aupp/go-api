package main

import (
	"net/http"
)

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"` // give it a nickname for json field
	}{
		Status: "ok",
	}
	responseWithJSON(w, http.StatusOK, response)
}
