package contracts

type UserCreatedEvent struct {
	UserID string
	Email  string
}

type UserUpdatedEvent struct {
	UserID string
	Email  string
	Name   string
}
