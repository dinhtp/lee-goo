package health

// HealthResponse is the HTTP response shape for the liveness probe.
type HealthResponse struct {
	Status string `json:"status"`
}
