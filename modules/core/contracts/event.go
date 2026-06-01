package contracts

// ModuleInstalledEvent is published when a module is successfully installed.
type ModuleInstalledEvent struct {
	Name    string
	Version string
}

// ModuleEnabledEvent is published when a module transitions to enabled.
type ModuleEnabledEvent struct {
	Name string
}

// ModuleDisabledEvent is published when a module transitions to disabled.
type ModuleDisabledEvent struct {
	Name string
}

// ModuleUninstalledEvent is published when a module is uninstalled.
type ModuleUninstalledEvent struct {
	Name string
}
