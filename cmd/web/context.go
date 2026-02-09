package main

import "lawbook/internal/models"

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

// templateData holds data passed to HTML templates
type templateData struct {
	CurrentYear     int
	Flash           string
	Form            interface{}
	IsAuthenticated bool
	CSRFToken       string
	User            *models.User
}
