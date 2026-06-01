package user

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

// UserResponse is the HTTP response shape — password hash is deliberately excluded.
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
