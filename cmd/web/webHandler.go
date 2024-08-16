package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/helper"
	// "github.com/julienschmidt/httprouter"
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
		Count: web.Form.Count,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) postCount(w http.ResponseWriter, _ *http.Request) {
	web.Form.Count++
	data := &countData{
		Count: web.Form.Count,
	}
	response := fmt.Sprintf(`<div id="count">%d</div>`, data.Count)
	helper.ResponseWithHyperMedia(w, http.StatusOK, response)
}

func (web *Web) getForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/form.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = ts.ExecuteTemplate(w, "base", 12)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) postForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.Form.FirstName = r.FormValue("firstname")
	web.Form.LastName = r.FormValue("lastname")

	response := fmt.Sprintf(`
		<ui>
			<li>First Name: %s</li>
		 <li>Last Name: %s</li>
		</ui>
		`, web.Form.FirstName, web.Form.LastName)
	helper.ResponseWithHyperMedia(w, http.StatusOK, response)
}
