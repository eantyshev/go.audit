package event

import (
	"context"

	audit "github.com/eantyshev/go.audit"
)

type Repository interface {
	InsertEvent(context.Context, audit.EventBase) (audit.ID, error)
	FindEvents(context.Context, audit.QueryParams) ([]audit.Event, error)

	Close()
}
