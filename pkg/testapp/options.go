package testapp

// WithConfig is a placeholder for future per-test config overrides.
// Tests currently rely on env vars for DATABASE_* and AUTH_JWT_SECRET.
func WithConfig() Option {
	return func(o *options) {}
}
