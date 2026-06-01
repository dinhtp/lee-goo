package module

// ModuleResponse is the HTTP response shape for module information.
type ModuleResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Path    string `json:"path"`
}
