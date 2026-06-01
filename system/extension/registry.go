package extension

import "sync"

// ExtensionRegistry holds named extension points for inter-module hooks.
// Full generic implementation arrives in Phase 6.
type ExtensionRegistry struct {
	mu     sync.Mutex
	points map[string][]any
}

// NewExtensionRegistry creates an empty ExtensionRegistry.
func NewExtensionRegistry() *ExtensionRegistry {
	return &ExtensionRegistry{points: make(map[string][]any)}
}

// Register appends a handler to the named extension point.
// The priority parameter is reserved for Phase 6 ordered dispatch.
func (r *ExtensionRegistry) Register(name string, priority int, handler any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.points[name] = append(r.points[name], handler)
}

// Resolve returns all handlers registered under the given name.
func (r *ExtensionRegistry) Resolve(name string) []any {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.points[name]
}
