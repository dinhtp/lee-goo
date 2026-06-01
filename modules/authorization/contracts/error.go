package contracts

import "errors"

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrUnauthorized       = errors.New("unauthorized")
)
