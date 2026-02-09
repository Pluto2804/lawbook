package main

import (
	"html/template"
	"path/filepath"
	"time"

	"lawbook/internal/models"
)

// templateData holds all the data needed by templates
type templateCache map[string]*template.Template

// newTemplateCache creates a cache of parsed templates
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all page templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a template set with custom functions
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// Template functions available in templates
var functions = template.FuncMap{
	"humanDate":   humanDate,
	"roleDisplay": roleDisplay,
}

// humanDate returns a nicely formatted string representation of a time.Time
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 at 15:04")
}

// roleDisplay returns a human-readable version of the role
func roleDisplay(role models.UserRole) string {
	switch role {
	case models.RoleStudent:
		return "Student"
	case models.RoleLawyer:
		return "Lawyer"
	case models.RoleRecruiter:
		return "Recruiter"
	default:
		return string(role)
	}
}
