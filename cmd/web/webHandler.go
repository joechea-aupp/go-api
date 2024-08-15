package web

import (
	"fmt"
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

type countData struct {
	Count int
}

func (web *Web) count(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/count.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := &countData{
		Count: web.Count,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) postCount(w http.ResponseWriter, _ *http.Request) {
	web.Count++
	data := &countData{
		Count: web.Count,
	}
	response := fmt.Sprintf(`<div id="count">%d</div>`, data.Count)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
