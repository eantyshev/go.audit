package handlers

import (
	"time"

	audit "github.com/eantyshev/go.audit"
)

type CreateEventRequest struct {
	Type     string                 `json:"type" binding:"required"`
	Consumer string                 `json:"consumer" binding:"required"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
}

type Event struct {
	Id        audit.ID  `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Type     string                 `json:"type" binding:"required"`
	Consumer string                 `json:"consumer" binding:"required"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
}

type QueryParams struct {
	Type, Consumer *string
	CreatedFrom    *time.Time `json:"created_from"`
	CreatedTo      *time.Time `json:"created_to"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type ListEventsResponse struct {
	Events []audit.Event `json:"events"`
}
