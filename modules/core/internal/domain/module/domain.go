package module

import "time"

// Status represents the lifecycle state of a module.
type Status string

const (
	StatusDiscovered  Status = "discovered"
	StatusInstalled   Status = "installed"
	StatusEnabled     Status = "enabled"
	StatusDisabled    Status = "disabled"
	StatusUpgrading   Status = "upgrading"
	StatusFailed      Status = "failed"
	StatusUninstalled Status = "uninstalled"
	StatusRemoved     Status = "removed"
)

// Module is the aggregate root for a managed module.
type Module struct {
	Name                  string
	Version               string
	Status                Status
	Path                  string
	Checksum              string
	InstalledAt           *time.Time
	EnabledAt             *time.Time
	DisabledAt            *time.Time
	UpgradedAt            *time.Time
	UninstalledAt         *time.Time
	RemovedFromCodebaseAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// Dependency describes a module dependency with its required/optional flag.
type Dependency struct {
	Name     string
	Required bool
}

// Manifest represents the parsed module.yaml descriptor.
type Manifest struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Status      string `yaml:"status"`
	Dependencies struct {
		Required []string `yaml:"required"`
		Optional []string `yaml:"optional"`
	} `yaml:"dependencies"`
	Provides struct {
		Services []string `yaml:"services"`
		Events   []string `yaml:"events"`
	} `yaml:"provides"`
	ExtensionPoints []string `yaml:"extension_points"`
	Migrations      struct {
		Path          string `yaml:"path"`
		Transactional bool   `yaml:"transactional"`
	} `yaml:"migrations"`
	Config struct {
		Prefix string `yaml:"prefix"`
	} `yaml:"config"`
}
