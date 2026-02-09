package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// serverError logs the error and sends a 500 Internal Server Error response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends a 404 Not Found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// renderer renders a template with the provided data
func (app *application) renderer(w http.ResponseWriter, req *http.Request, page string, status int, data *templateData) {
	ts, ok := app.tempCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

// routerWrap creates a handler for router's NotFound and MethodNotAllowed
func (app *application) routerWrap() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		app.notFound(w)
	}
}

// decodePostForm decodes POST form data into a destination struct
func (app *application) decodePostForm(req *http.Request, dst interface{}) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, req.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

// isAuthenticated checks if the current request is from an authenticated user
func (app *application) isAuthenticated(req *http.Request) bool {
	isAuthenticated, ok := req.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

// newTemplateData creates a new templateData struct with default values
func (app *application) newTemplateData(req *http.Request) *templateData {
	data := &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(req.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(req),
		CSRFToken:       nosurf.Token(req),
	}

	// Add user info if authenticated
	if data.IsAuthenticated {
		userID := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")
		user, err := app.models.Users.Get(userID)
		if err == nil {
			data.User = user
		}
	}

	return data
}
