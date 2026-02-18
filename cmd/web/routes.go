package main

import (
	"net/http"

	"lawbook/internal/models"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = app.routerWrap()
	router.MethodNotAllowed = app.routerWrap()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Base middleware chain (security, logging, panic recovery)
	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)

	// Dynamic middleware (with sessions, CSRF, and auth check)
	dynamic := standard.Append(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)

	// Protected routes (requires authentication)
	protected := dynamic.Append(app.requireAuthentication)

	// Role-specific middleware chains
	studentOnly := protected.Append(app.requireRole(models.RoleStudent))
	lawyerOnly := protected.Append(app.requireRole(models.RoleLawyer))
	recruiterOnly := protected.Append(app.requireRole(models.RoleRecruiter))

	// Lawyers and students can access moot court
	mootCourtAccess := protected.Append(app.requireAnyRole(models.RoleStudent, models.RoleLawyer))

	// ==================== PUBLIC ROUTES ====================
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	// Add under PUBLIC ROUTES
	router.Handler(http.MethodGet, "/api/user/me", dynamic.ThenFunc(app.apiUserMe))

	// Authentication routes
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// ==================== PROTECTED ROUTES ====================
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))
	router.Handler(http.MethodGet, "/user/account", protected.ThenFunc(app.accountView))

	// ==================== STUDENT ROUTES ====================
	router.Handler(http.MethodGet, "/student/dashboard", studentOnly.ThenFunc(app.studentDashboard))

	// ==================== LAWYER ROUTES ====================
	router.Handler(http.MethodGet, "/lawyer/dashboard", lawyerOnly.ThenFunc(app.lawyerDashboard))

	// ==================== RECRUITER ROUTES ====================
	router.Handler(http.MethodGet, "/recruiter/dashboard", recruiterOnly.ThenFunc(app.recruiterDashboard))

	// ==================== MOOT COURT ROUTES (Students & Lawyers) ====================
	router.Handler(http.MethodGet, "/moot/setup", mootCourtAccess.ThenFunc(app.mootCourtSetup))
	router.Handler(http.MethodGet, "/moot/session", mootCourtAccess.ThenFunc(app.mootCourtSession))

	return dynamic.Then(router)
}
