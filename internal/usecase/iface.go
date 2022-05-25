package usecase

import (
	"go.audit/internal/entity"
)

type Iface interface {
	AddEvent(event entity.Event) error
	FindEvents(params entity.QueryParams) ([]entity.Event, error)
}
