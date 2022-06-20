package repository

import (
	"context"
	"strconv"
	"sync"
	"time"

	audit "github.com/eantyshev/go.audit"
)

type Matcher struct {
	params audit.QueryParams
}

func (m *Matcher) Matches(event audit.Event) bool {
	if m.params.Consumer != nil && *m.params.Consumer != event.Consumer {
		return false
	}
	if m.params.Type != nil && *m.params.Type != event.Type {
		return false
	}
	return true
}

type MemRepo struct {
	mx    sync.RWMutex
	store []audit.Event
}

func (r *MemRepo) InsertEvent(_ context.Context, event_base audit.EventBase) (audit.ID, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var idx int = len(r.store)
	event := audit.Event{EventBase: event_base}
	event.Id = audit.ID(strconv.Itoa(idx))
	event.CreatedAt = time.Now()
	r.store = append(r.store, event)
	return event.Id, nil
}

func (r *MemRepo) FindEvents(_ context.Context, params audit.QueryParams) (lst []audit.Event, err error) {
	r.mx.RLock()
	defer r.mx.Unlock()

	lst = make([]audit.Event, 0, len(r.store))
	matcher := Matcher{params}
	for _, event := range r.store {
		if matcher.Matches(event) {
			lst = append(lst, event)
		}
	}
	return lst, nil
}

func (r *MemRepo) Close() {
	// wait for pending requests
	r.mx.Lock()
	defer r.mx.Unlock()
}
