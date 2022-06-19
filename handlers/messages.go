package handlers

import (
	"time"

	"go.audit/internal/entity"
)

type CreateEventRequest struct {
	Type     string                 `json:"type" binding:"required"`
	Consumer string                 `json:"consumer" binding:"required"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
}

type Event struct {
	Id        ID        `json:"id"`
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
	Events []entity.Event `json:"events"`
}
