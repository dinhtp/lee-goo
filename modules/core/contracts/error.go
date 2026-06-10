package contracts

import "errors"

var (
	ErrModuleNotFound      = errors.New("module not found")
	ErrModuleAlreadyExists = errors.New("module already exists")
	ErrCircularDependency  = errors.New("circular dependency detected")
	ErrModuleHasDependents = errors.New("module has dependents that require it")
	ErrProtectedModule     = errors.New("module is protected and cannot be modified")
	ErrInvalidModulePath   = errors.New("invalid module path")
	ErrModuleNotInstalled  = errors.New("module is not installed")
)
