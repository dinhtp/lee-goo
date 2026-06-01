package contracts

// LoginSucceededEvent is published when a user authenticates successfully.
type LoginSucceededEvent struct {
	UserID string
	Email  string
}

// LoginFailedEvent is published when authentication fails (wrong credentials).
type LoginFailedEvent struct {
	Email string
}
