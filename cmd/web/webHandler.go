package web

import (
	"fmt"
	"net/http"

	"github.com/joechea-aupp/go-api/cmd/helper"
	"github.com/joechea-aupp/go-api/ui"
	"github.com/julienschmidt/httprouter"
)

func (web *Web) render(w http.ResponseWriter, status int, page string, data *ui.TemplateData) {
	ts, ok := web.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) user(w http.ResponseWriter, r *http.Request) {
	web.render(w, http.StatusOK, "user.tmpl.html", nil)
}

func (web *Web) count(w http.ResponseWriter, r *http.Request) {
	web.render(w, http.StatusOK, "count.tmpl.html", web.templateData)
}

func (web *Web) postCount(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	if params.ByName("mode") == "increment" {
		web.templateData.Count++
	} else if params.ByName("mode") == "decrement" {
		web.templateData.Count--
	} else {
		helper.ResponseWithError(w, http.StatusBadRequest, "invalid mode")
		return
	}

	response := fmt.Sprintf(`
  <h1 class="text-9xl" id="count">%d</h1>
`, web.templateData.Count)

	helper.ResponseWithHyperMedia(w, http.StatusOK, response)
}

func (web *Web) getForm(w http.ResponseWriter, r *http.Request) {
	web.render(w, http.StatusOK, "form.tmpl.html", web.templateData)
}

// func (web *Web) postForm(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	web.Form.FirstName = r.FormValue("firstname")
// 	web.Form.LastName = r.FormValue("lastname")
//
// 	response := fmt.Sprintf(`
// 		<ui>
// 			<li>First Name: %s</li>
// 		 <li>Last Name: %s</li>
// 		</ui>
// 		`, web.Form.FirstName, web.Form.LastName)
// 	helper.ResponseWithHyperMedia(w, http.StatusOK, response)
// }
