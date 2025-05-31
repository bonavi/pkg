package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"

	"pkg/log"
)

func (m *Monitor) IncFailedEvent(_ context.Context, evt *event.CommandFailedEvent) {
	if globalMetric.mongoCommandFailedMetric == nil {
		log.Error("mongoCommandFailedMetric prometheus metric not initialized")
		return
	}

	m.mu.Lock()
	cmd := m.commands[evt.RequestID]
	delete(m.commands, evt.RequestID)
	m.mu.Unlock()

	globalMetric.mongoCommandFailedMetric.WithLabelValues(m.namespace, cmd.database, cmd.collection, evt.CommandName).Inc()
}
