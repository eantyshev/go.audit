package repository

import (
	"context"

	"go.audit/internal/entity"
)

type RepoIface interface {
	InsertEvent(context.Context, entity.Event) (entity.ID, error)
	FindEvents(context.Context, entity.QueryParams) ([]entity.Event, error)

	Close()
}
