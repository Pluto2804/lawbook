package main

import (
	"errors"
	"net/http"

	"lawbook/internal/models"
	"lawbook/internal/validator"
)

// ==================== HOME & PUBLIC PAGES ====================

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "home.tmpl.html", http.StatusOK, data)
}

func (app *application) about(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "about.tmpl.html", http.StatusOK, data)
}

// ==================== USER SIGNUP ====================

type userSignupForm struct {
	Name                string          `form:"name"`
	Email               string          `form:"email"`
	Password            string          `form:"password"`
	Role                models.UserRole `form:"role"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = userSignupForm{}
	app.renderer(w, req, "signup.tmpl.html", http.StatusOK, data)
}

func (app *application) userSignupPost(w http.ResponseWriter, req *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validation
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// Validate role
	validRole := form.Role == models.RoleStudent ||
		form.Role == models.RoleLawyer ||
		form.Role == models.RoleRecruiter
	form.CheckField(validRole, "role", "Please select a valid role")

	if !form.Valid() {
		data := app.newTemplateData(req)
		data.Form = form
		app.renderer(w, req, "signup.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}

	// Insert user
	_, err = app.models.Users.Insert(form.Name, form.Email, form.Password, form.Role)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldErrors("email", "Email address is already in use")
			data := app.newTemplateData(req)
			data.Form = form
			app.renderer(w, req, "signup.tmpl.html", http.StatusUnprocessableEntity, data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(req.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, req, "/user/login", http.StatusSeeOther)
}

// ==================== USER LOGIN ====================

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = userLoginForm{}
	app.renderer(w, req, "login.tmpl.html", http.StatusOK, data)
}

func (app *application) userLoginPost(w http.ResponseWriter, req *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(req)
		data.Form = form
		app.renderer(w, req, "login.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}

	id, err := app.models.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(req)
			data.Form = form
			app.renderer(w, req, "login.tmpl.html", http.StatusUnprocessableEntity, data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Renew session token
	err = app.sessionManager.RenewToken(req.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Store user ID in session
	app.sessionManager.Put(req.Context(), "authenticatedUserId", id)

	// Get user to determine role-based redirect
	user, err := app.models.Users.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect based on role
	switch user.Role {
	case models.RoleStudent:
		http.Redirect(w, req, "/student/dashboard", http.StatusSeeOther)
	case models.RoleLawyer:
		http.Redirect(w, req, "/lawyer/dashboard", http.StatusSeeOther)
	case models.RoleRecruiter:
		http.Redirect(w, req, "/recruiter/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

// ==================== USER LOGOUT ====================

func (app *application) userLogout(w http.ResponseWriter, req *http.Request) {
	err := app.sessionManager.RenewToken(req.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(req.Context(), "authenticatedUserId")
	app.sessionManager.Put(req.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// ==================== ACCOUNT MANAGEMENT ====================

func (app *application) accountView(w http.ResponseWriter, req *http.Request) {
	userID := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")

	user, err := app.models.Users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(req)
	data.User = user
	app.renderer(w, req, "account.tmpl.html", http.StatusOK, data)
}

// ==================== ROLE-SPECIFIC DASHBOARDS ====================

// Student Dashboard
func (app *application) studentDashboard(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "student-dashboard.tmpl.html", http.StatusOK, data)
}

// Lawyer Dashboard
func (app *application) lawyerDashboard(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "lawyer-dashboard.tmpl.html", http.StatusOK, data)
}

// Recruiter Dashboard
func (app *application) recruiterDashboard(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "recruiter-dashboard.tmpl.html", http.StatusOK, data)
}

// ==================== MOOT COURT SIMULATOR ====================

func (app *application) mootCourtSetup(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "moot-setup.tmpl.html", http.StatusOK, data)
}

func (app *application) mootCourtSession(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "moot-session.tmpl.html", http.StatusOK, data)
}
