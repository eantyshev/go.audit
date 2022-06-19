package services

import (
	"context"
	"sort"

	"go.audit/entity"
)

type Usecase struct {
	Repo repository.RepoIface
}

var _ Iface = &Usecase{}

func (u *Usecase) AddEvent(event entity.Event) error {
	_, err := u.Repo.InsertEvent(context.TODO(), event)
	return err
}

func (u *Usecase) FindEvents(params entity.QueryParams) ([]entity.Event, error) {
	events, err := u.Repo.FindEvents(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	return events, nil
}
