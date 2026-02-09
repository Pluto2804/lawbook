package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleStudent   UserRole = "student"
	RoleLawyer    UserRole = "lawyer"
	RoleRecruiter UserRole = "recruiter"
)

// User represents a user in the system
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Role           UserRole
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsActive       bool
	EmailVerified  bool
}

// UserModel wraps a database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new user to the database
func (m *UserModel) Insert(name, email, password string, role UserRole) (int, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, role)
		VALUES (?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, email, hashedPassword, role)
	if err != nil {
		// Check for duplicate email error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Authenticate verifies a user's email and password
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	var isActive bool

	stmt := `SELECT id, hashed_password, is_active FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Check if user account is active
	if !isActive {
		return 0, ErrInactiveAccount
	}

	// Compare the hashed password with the plain-text password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get retrieves a user by their ID
func (m *UserModel) Get(id int) (*User, error) {
	stmt := `SELECT id, name, email, role, created_at, updated_at, is_active, email_verified
		FROM users WHERE id = ?`

	var user User

	err := m.DB.QueryRow(stmt, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.EmailVerified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

// Exists checks if a user with a given ID exists
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

// UpdatePassword changes a user's password
func (m *UserModel) UpdatePassword(id int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET hashed_password = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err = m.DB.Exec(stmt, hashedPassword, id)
	return err
}

// VerifyEmail marks a user's email as verified
func (m *UserModel) VerifyEmail(id int) error {
	stmt := `UPDATE users SET email_verified = TRUE, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	return err
}

// DeactivateUser sets a user's account to inactive
func (m *UserModel) DeactivateUser(id int) error {
	stmt := `UPDATE users SET is_active = FALSE, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	return err
}

// GetByRole retrieves all users with a specific role (useful for admin functions)
func (m *UserModel) GetByRole(role UserRole, limit, offset int) ([]*User, error) {
	stmt := `SELECT id, name, email, role, created_at, updated_at, is_active, email_verified
		FROM users WHERE role = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(stmt, role, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsActive,
			&user.EmailVerified,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
