package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/joechea-aupp/go-api/cmd/helper"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/joechea-aupp/go-api/internal/validator"
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

func (web *Web) users(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()

	startParam := queryParam.Get("start")
	if startParam == "" {
		startParam = "0"
	}

	limitParam := queryParam.Get("limit")
	if limitParam == "" {
		limitParam = "2"
	}

	start, err := strconv.ParseInt(startParam, 10, 64)
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := web.User.GetUsers(start, limit)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.templateData.Users = users

	// get total number of users
	userCount, err := web.User.TotalUsers()
	if err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	web.templateData.Form = struct {
		TotalUsers int64
		Start      int
	}{
		TotalUsers: userCount,
		Start:      int(start),
	}

	web.templateData.Flash = web.sessionManager.PopString(r.Context(), "flash")
	web.render(w, http.StatusOK, "users.tmpl.html", web.templateData)
}

func (web *Web) user(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	username := params.ByName("username")

	user, err := web.User.GetUser(username)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := &db.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	web.templateData.User = data
	web.render(w, http.StatusOK, "user.tmpl.html", web.templateData)
}

func (web *Web) updateUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := r.ParseForm()
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	username := r.Form.Get("floating_username")
	email := r.Form.Get("floating_email")

	user := db.User{
		Username: username,
		Email:    email,
	}

	err = web.User.UpdateUser(id, user)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.sessionManager.Put(r.Context(), "flash", "User updated successfully")

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (web *Web) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := web.User.DelUser(id)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.sessionManager.Put(r.Context(), "flash", "User deleted successfully")

	http.Redirect(w, r, "/users", http.StatusSeeOther)
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

func (web *Web) postForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	username := r.Form.Get("floating_username")
	email := r.Form.Get("floating_email")
	password := r.Form.Get("floating_password")
	repeatPassword := r.Form.Get("repeat_password")

	if password != repeatPassword {
		helper.ResponseWithError(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	user := db.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	err = web.User.CreateUser(user)
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.sessionManager.Put(r.Context(), "flash", "User created successfully")

	// redirect user to /user with status code of 303
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (web *Web) formValidator(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	field := params.ByName("field")

	err := r.ParseForm()
	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := ""

	submitBtn := fmt.Sprintf(`
      <button
        disabled
        type="submit"
        hx-swap-oob="outerHTML"
        id="submit-button"
        class="py-2.5 px-5 me-2 mb-2 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-200 dark:bg-gray-500 dark:text-gray-400 dark:border-gray-300"
      >
      Submit
    </button>
		`)

	if _, ok := r.Form[field]; !ok {
		helper.ResponseWithError(w, http.StatusBadRequest, "field does not exist")
		return
	}

	if field == "floating_username" {
		username := r.Form.Get("floating_username")
		if !validator.MinChars(username, 3) {
			message = "Username must be at least 3 characters long"
			errorResponse := fmt.Sprintf(`
		<p class="validate_field mt-2 text-sm text-red-600 dark:text-red-500"><span class="font-medium">Oops!</span>%s</p>
	`, message)

			errorResponse += submitBtn

			helper.ResponseWithHyperMedia(w, http.StatusOK, errorResponse)
			return
		}
	}

	if field == "floating_email" {
		email := r.Form.Get("floating_email")
		if !validator.Matches(email, validator.EmailRX) {
			message = "Invalid email address"
			errorResponse := fmt.Sprintf(`
		<p class="validate_field mt-2 text-sm text-red-600 dark:text-red-500"><span class="font-medium">Oops!</span>%s</p>
	`, message)

			errorResponse += submitBtn

			helper.ResponseWithHyperMedia(w, http.StatusOK, errorResponse)
			return
		}
	}

	submitBtn = `
	<button
      type="submit"
      hx-swap-oob="true"
      id="submit-button"
      class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
    >
      Submit
    </button>
	`

	message = fmt.Sprintf(`
				<p class="validate_field"></p>
			`)

	message += submitBtn

	helper.ResponseWithHyperMedia(w, http.StatusOK, message)
}
