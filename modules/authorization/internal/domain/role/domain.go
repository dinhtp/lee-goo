package role

import "time"

// Role represents an authorization role that groups a set of permissions.
type Role struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

// Permission is a discrete action-resource pair (e.g. "create" on "orders").
type Permission struct {
	ID       string
	Action   string // e.g. "create", "read", "update", "delete"
	Resource string // e.g. "users", "orders"
}

// UserRole records the assignment of a Role to a User.
type UserRole struct {
	UserID     string
	RoleID     string
	AssignedAt time.Time
}
