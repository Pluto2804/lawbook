package models

import "database/sql"

// Models wraps all the model types
type Models struct {
	Users    *UserModel
	Sessions *SessionModel
}

// NewModels returns a Models struct containing initialized model types
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:    &UserModel{DB: db},
		Sessions: &SessionModel{DB: db},
	}
}
