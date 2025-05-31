package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"

	"pkg/log"
)

func (m *Monitor) IncSucceededEvent(_ context.Context, evt *event.CommandSucceededEvent) {
	if globalMetric.mongoCommandSucceededMetric == nil {
		log.Error("mongoCommandSucceededMetric prometheus metric not initialized")
		return
	}

	m.mu.Lock()
	cmd := m.commands[evt.RequestID]
	delete(m.commands, evt.RequestID)
	m.mu.Unlock()

	globalMetric.mongoCommandSucceededMetric.WithLabelValues(
		m.namespace, cmd.database, cmd.collection, evt.CommandName,
	).Observe(evt.Duration.Seconds())
}
