package contracts

import "context"

// IdentityProvider is implemented by authentication middleware to expose
// the current authenticated user to the authorization layer.
type IdentityProvider interface {
	CurrentUser(ctx context.Context) (*UserIdentity, error)
}

// UserIdentity carries the resolved identity from a token or session.
type UserIdentity struct {
	UserID string
	Email  string
	Roles  []string
}
