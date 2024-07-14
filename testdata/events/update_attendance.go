package events

import "github.com/jpicht/polyjson"

type UpdateAttendance struct {
	polyjson.Implements[Event]

	DateID  int             `json:"date_id"`
	Value   AttendanceState `json:"value"`
	Comment string          `json:"comment,omitempty"`
}

type AttendanceState string

const (
	WillAttend    = AttendanceState("yes")
	WillNotAttend = AttendanceState("no")
	MayAttend     = AttendanceState("maybe")
)
