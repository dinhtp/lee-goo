package eventbus

import "context"

// Handler is a function that processes an event payload.
type Handler func(ctx context.Context, payload any) error

// EventBus defines the publish/subscribe contract for inter-module events.
type EventBus interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(topic string, handler Handler)
}

// NoopEventBus is a no-op stub — replaced by the real implementation in Phase 6.
type NoopEventBus struct{}

// compile-time interface check
var _ EventBus = (*NoopEventBus)(nil)

func (n *NoopEventBus) Publish(_ context.Context, _ string, _ any) error { return nil }
func (n *NoopEventBus) Subscribe(_ string, _ Handler)                    {}

// NewNoopEventBus returns the no-op EventBus implementation.
func NewNoopEventBus() EventBus { return &NoopEventBus{} }
