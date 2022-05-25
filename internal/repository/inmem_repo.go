package repository

import (
	"context"
	"strconv"
	"time"

	"go.audit/internal/entity"
)

type Matcher struct {
	params entity.QueryParams
}

func (m *Matcher) Matches(event entity.Event) bool {
	if m.params.Consumer != nil && *m.params.Consumer != event.Consumer {
		return false
	}
	if m.params.Type != nil && *m.params.Type != event.Type {
		return false
	}
	return true
}

type MemRepo struct {
	store []entity.Event
}

func (r *MemRepo) InsertEvent(_ context.Context, event entity.Event) (entity.ID, error) {
	var idx int = len(r.store)
	event.Id = entity.ID(strconv.Itoa(idx))
	event.CreatedAt = time.Now()
	r.store = append(r.store, event)
	return event.Id, nil
}

func (r *MemRepo) FindEvents(_ context.Context, params entity.QueryParams) (lst []entity.Event, err error) {
	lst = make([]entity.Event, 0, len(r.store))
	matcher := Matcher{params}
	for _, event := range r.store {
		if matcher.Matches(event) {
			lst = append(lst, event)
		}
	}
	return lst, nil
}

func (r *MemRepo) Close() {}
