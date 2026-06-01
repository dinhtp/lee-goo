package user

import "time"

// User is the internal domain entity. Never expose this directly across modules.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
