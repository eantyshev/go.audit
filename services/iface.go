package services

import audit "github.com/eantyshev/go.audit"

type EventSvc interface {
	AddEvent(event audit.EventBase) error
	FindEvents(params audit.QueryParams) ([]audit.Event, error)
}
