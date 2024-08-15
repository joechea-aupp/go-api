package web

import (
	"html/template"
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/helper"
)

func (web *Web) user(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) count(w http.ResponseWriter, r *http.Request) {
}
