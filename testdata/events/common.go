package events

import (
	"time"

	"github.com/jpicht/polyjson"
)

//go:generate go run github.com/jpicht/polyjson/cmd/gopolyjson -F .

type Common struct {
	polyjson.Common[Event]

	Time         time.Time `json:"time"`
	UserID       int       `json:"user_id,omitempty"`
	ActualUserID int       `json:"actual_user_id,omitempty"`
}
