package config

// ModuleConfig holds configuration for the module manager.
type ModuleConfig struct {
	ProtectedModules   []string `mapstructure:"protected_modules"`
	AllowSourceRemoval bool     `mapstructure:"allow_source_removal"`
}

// DefaultConfig returns the default module manager configuration.
func DefaultConfig() *ModuleConfig {
	return &ModuleConfig{
		ProtectedModules:   []string{"module", "user", "authentication", "authorization"},
		AllowSourceRemoval: false,
	}
}
