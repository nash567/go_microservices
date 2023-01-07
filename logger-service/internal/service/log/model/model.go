package model

import "time"

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name,omitempty" json:"name"`
	Data      string    `bson:"data,omitempty" json:"data"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}
type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
