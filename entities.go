package go.audit

import "time"

type ID string

type EventBase struct {
	Type     string
	Consumer string
	Payload  map[string]interface{}
}

type Event struct {
	Id        ID
	CreatedAt time.Time
	EventBase
}

type QueryParams struct {
	Type, Consumer string
	CreatedFrom    time.Time
	CreatedTo      time.Time
}
