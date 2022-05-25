package api

import "go.audit/internal/entity"

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type ListEventsResponse struct {
	Events []entity.Event `json:"events"`
}
