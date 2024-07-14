package events

import (
	"time"

	_ "go.uber.org/mock/gomock"

	"github.com/jpicht/polyjson"
)

//go:generate go run github.com/jpicht/polyjson/cmd/gopolyjson -F .
//go:generate go run go.uber.org/mock/mockgen -package events -destination eventvisitor_mock.go . EventVisitor

type Common struct {
	polyjson.Common[Event]

	Time         time.Time `json:"time"`
	UserID       int       `json:"user_id,omitempty"`
	ActualUserID int       `json:"actual_user_id,omitempty"`
}
