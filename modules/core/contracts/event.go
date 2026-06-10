package contracts

// ModuleInstalledEvent is published when a module is successfully installed.
type ModuleInstalledEvent struct {
	Name    string
	Version string
}

// ModuleUninstalledEvent is published when a module is uninstalled.
type ModuleUninstalledEvent struct {
	Name string
}
