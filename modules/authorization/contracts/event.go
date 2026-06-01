package contracts

// RoleAssignedEvent is published when a role is assigned to a user.
type RoleAssignedEvent struct {
	UserID string
	RoleID string
}

// RoleRemovedEvent is published when a role is removed from a user.
type RoleRemovedEvent struct {
	UserID string
	RoleID string
}
