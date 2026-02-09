package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"time"
)

// Session represents a user session
type Session struct {
	Token  string
	UserID int
	Expiry time.Time
}

// SessionModel wraps a database connection pool
type SessionModel struct {
	DB *sql.DB
}

// Insert creates a new session for a user
func (m *SessionModel) Insert(userID int) (string, error) {
	// Generate a random session token
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	// Set session expiry to 12 hours from now
	expiry := time.Now().Add(12 * time.Hour)

	stmt := `INSERT INTO sessions (token, user_id, expiry)
		VALUES (?, ?, ?)`

	_, err = m.DB.Exec(stmt, token, userID, expiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Get retrieves the user ID associated with a session token
func (m *SessionModel) Get(token string) (int, error) {
	var userID int
	var expiry time.Time

	stmt := `SELECT user_id, expiry FROM sessions WHERE token = ?`

	err := m.DB.QueryRow(stmt, token).Scan(&userID, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNoRecord
		}
		return 0, err
	}

	// Check if the session has expired
	if time.Now().After(expiry) {
		return 0, ErrExpiredSession
	}

	return userID, nil
}

// Delete removes a session from the database
func (m *SessionModel) Delete(token string) error {
	stmt := `DELETE FROM sessions WHERE token = ?`

	_, err := m.DB.Exec(stmt, token)
	return err
}

// DeleteAllForUser removes all sessions for a specific user
func (m *SessionModel) DeleteAllForUser(userID int) error {
	stmt := `DELETE FROM sessions WHERE user_id = ?`

	_, err := m.DB.Exec(stmt, userID)
	return err
}

// CleanupExpired removes all expired sessions from the database
func (m *SessionModel) CleanupExpired() error {
	stmt := `DELETE FROM sessions WHERE expiry < UTC_TIMESTAMP()`

	_, err := m.DB.Exec(stmt)
	return err
}

// generateSessionToken creates a cryptographically secure random token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode to base32 for URL-safe token
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return token, nil
}
