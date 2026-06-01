package role

// CreateRoleRequest is the request body for POST /roles.
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AssignRoleRequest is the request body for POST /users/:id/roles.
type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

// RoleResponse is the JSON representation returned by role endpoints.
type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
