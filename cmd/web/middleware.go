package main

import (
	"context"
	"fmt"
	"net/http"

	"lawbook/internal/models"

	"github.com/justinas/nosurf"
)

// secureHeaders adds security headers to every response
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://lawbookv2.vercel.app")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, req)
	})
}

// logRequest logs each HTTP request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(w, req)
	})
}

// recoverPanic recovers from panics and sends a 500 Internal Server Error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

// requireAuthentication checks if the user is authenticated
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !app.isAuthenticated(req) {
			app.sessionManager.Put(req.Context(), "redirectPathAfterLogin", req.URL.Path)
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
			return
		}

		// Tell browser not to cache pages for authenticated users
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, req)
	})
}

// noSurf provides CSRF protection
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "CSRF token validation failed: "+nosurf.Reason(r).Error(), http.StatusBadRequest)
	}))

	return csrfHandler
}

// authenticate checks if a user is authenticated and adds user info to context
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")
		if id == 0 {
			next.ServeHTTP(w, req)
			return
		}

		exists, err := app.models.Users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(req.Context(), isAuthenticatedContextKey, true)
			req = req.WithContext(ctx)
		}

		next.ServeHTTP(w, req)
	})
}

// requireRole creates middleware that checks if user has a specific role
func (app *application) requireRole(role models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			userID := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")

			user, err := app.models.Users.Get(userID)
			if err != nil {
				app.serverError(w, err)
				return
			}

			if user.Role != role {
				app.sessionManager.Put(req.Context(), "flash", "You don't have permission to access this page")
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// requireAnyRole creates middleware that checks if user has any of the specified roles
func (app *application) requireAnyRole(roles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			userID := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")

			user, err := app.models.Users.Get(userID)
			if err != nil {
				app.serverError(w, err)
				return
			}

			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				app.sessionManager.Put(req.Context(), "flash", "You don't have permission to access this page")
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
