package eventbus

import (
	"context"
	"sync"

	"github.com/dinhtp/lee-goo/system/logger"
)

// localEventBus is an in-process publish/subscribe implementation of EventBus.
// All handlers for a topic are invoked synchronously in Publish order.
// Handler failures are logged but do not abort delivery to remaining handlers.
type localEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	logger   *logger.Logger
}

// compile-time interface check
var _ EventBus = (*localEventBus)(nil)

// NewLocalEventBus constructs an in-process EventBus backed by a handler map.
func NewLocalEventBus(logger *logger.Logger) EventBus {
	return &localEventBus{
		handlers: make(map[string][]Handler),
		logger:   logger,
	}
}

// Subscribe registers h to receive events on topic.
func (b *localEventBus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}

// Publish dispatches payload to all handlers subscribed to topic.
// A handler error is logged but delivery continues to remaining handlers.
func (b *localEventBus) Publish(ctx context.Context, topic string, payload any) error {
	b.mu.RLock()
	// Copy slice under read lock to avoid holding it during handler execution.
	handlers := make([]Handler, len(b.handlers[topic]))
	copy(handlers, b.handlers[topic])
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, payload); err != nil {
			b.logger.Sugar().Errorw("event handler error", "topic", topic, "error", err)
			// Continue to other handlers — one failure does not abort publish.
		}
	}
	return nil
}
