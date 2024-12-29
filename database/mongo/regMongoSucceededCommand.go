package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"

	"pkg/log"
)

func (m *Monitor) IncSucceededEvent(ctx context.Context, event *event.CommandSucceededEvent) {
	if globalMetric.mongoCommandSucceededMetric == nil {
		log.Error(ctx, "mongoCommandSucceededMetric prometheus metric not initialized")
		return
	}

	m.mu.Lock()
	cmd := m.commands[event.RequestID]
	delete(m.commands, event.RequestID)
	m.mu.Unlock()

	globalMetric.mongoCommandSucceededMetric.WithLabelValues(
		m.namespace, cmd.database, cmd.collection, event.CommandName,
	).Observe(event.Duration.Seconds())
}
