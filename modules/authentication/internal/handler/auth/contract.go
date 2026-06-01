package auth

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest is the JSON body for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse is the JSON response body for successful authentication operations.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}
