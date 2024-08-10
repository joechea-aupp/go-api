package helper

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	errorMsg := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	ResponseWithJSON(w, code, errorMsg)
}
