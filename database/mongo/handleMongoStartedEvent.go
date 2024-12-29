package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
)

func (m *Monitor) HandleStartedEvent(_ context.Context, evt *event.CommandStartedEvent) {
	collectionRaw := evt.Command.Lookup(evt.CommandName)
	collection, _ := collectionRaw.StringValueOK()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.commands[evt.RequestID] = command{
		database:   evt.DatabaseName,
		collection: collection,
		raw:        evt.Command,
	}
}
