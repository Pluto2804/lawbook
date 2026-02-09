package models

import "errors"

var (
	// ErrNoRecord is returned when a record doesn't exist in the database
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials is returned when a user provides invalid login credentials
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail is returned when trying to create a user with an email that already exists
	ErrDuplicateEmail = errors.New("models: duplicate email")

	// ErrInactiveAccount is returned when a user's account is deactivated
	ErrInactiveAccount = errors.New("models: account is inactive")

	// ErrExpiredSession is returned when a session has expired
	ErrExpiredSession = errors.New("models: session has expired")
)
