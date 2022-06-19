package event

import (
	"context"

	"go.audit/entity"
)

type Repository interface {
	InsertEvent(context.Context, entity.Event) (entity.ID, error)
	FindEvents(context.Context, entity.QueryParams) ([]entity.Event, error)

	Close()
}
