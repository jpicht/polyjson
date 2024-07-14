package typefield

import (
	"time"

	_ "go.uber.org/mock/gomock"

	"github.com/jpicht/polyjson"
)

//go:generate go run github.com/jpicht/polyjson/cmd/gopolyjson -F .
//go:generate go run go.uber.org/mock/mockgen -package typefield -destination eventvisitor_mock.go . EventVisitor

type EventType = polyjson.TypeID[Event]

type Common struct {
	polyjson.Common[Event]

	Time         time.Time `json:"time"`
	Type         EventType `json:"type"`
	UserID       int       `json:"user_id,omitempty"`
	ActualUserID int       `json:"actual_user_id,omitempty"`
}
