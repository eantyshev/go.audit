package entity

import "time"

type ID string

type Event struct {
	Id        ID                     `json:"id" bson:"_id,omitempty"`
	Type      string                 `json:"type" binding:"required" bson:"type"`
	Consumer  string                 `json:"consumer" binding:"required" bson:"consumer"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	Payload   map[string]interface{} `json:"payload,omitempty" bson:"payload,inline,omitempty"`
}

type QueryParams struct {
	Type, Consumer *string
	CreatedFrom    *time.Time `json:"created_from"`
	CreatedTo      *time.Time `json:"created_to"`
}
