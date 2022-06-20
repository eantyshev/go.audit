package services

import (
	"context"
	"sort"

	audit "github.com/eantyshev/go.audit"
	"github.com/eantyshev/go.audit/domain/event"
)

type EventSvcImpl struct {
	Repo event.Repository
}

var _ EventSvc = &EventSvcImpl{}

func (u *EventSvcImpl) AddEvent(event audit.EventBase) error {
	_, err := u.Repo.InsertEvent(context.TODO(), event)
	return err
}

func (u *EventSvcImpl) FindEvents(params audit.QueryParams) ([]audit.Event, error) {
	events, err := u.Repo.FindEvents(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	return events, nil
}
