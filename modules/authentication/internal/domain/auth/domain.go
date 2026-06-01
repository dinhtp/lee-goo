package auth

// Claims holds the decoded JWT claims used internally within the authentication domain.
type Claims struct {
	UserID string
	Email  string
	Type   string // "access" | "refresh"
}

// TokenPair is the internal domain representation of an issued token pair.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}
