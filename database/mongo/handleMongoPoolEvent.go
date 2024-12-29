package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"

	"pkg/log"
)

func HandlePoolEvent(evt *event.PoolEvent) {
	if globalMetric.mongoPoolEventsMetric == nil {
		log.Error(context.Background(), "mongoPoolEventsMetric prometheus metric not initialized")
		return
	}

	switch evt.Type {
	case event.PoolCreated:
		poolEvent("connection_pool_created")
	case event.PoolCleared:
		poolEvent("connection_pool_cleared")
	case event.PoolClosedEvent:
		poolEvent("connection_pool_closed")
	case event.ConnectionCreated:
		poolEvent("connection_created")
	case event.ConnectionReady:
		poolEvent("connection_ready")
	case event.ConnectionClosed:
		poolEvent("connection_closed")
	case event.GetStarted:
		poolEvent("connection_started")
	case event.GetFailed:
		poolEvent("connection_failed")
	case event.GetSucceeded:
		poolEvent("connection_succeeded")
	case event.ConnectionReturned:
		poolEvent("connection_returned")
	}
}

func poolEvent(event string) {
	globalMetric.mongoPoolEventsMetric.WithLabelValues(globalMetric.namespace, event).Inc()
}
