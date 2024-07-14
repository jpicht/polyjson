package polyjson_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/jpicht/polyjson/testdata/events"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:embed testdata/events/example.json
var jsondata []byte

func TestEvents(t *testing.T) {
	var e events.EventSlice

	require.NoError(t, json.Unmarshal(jsondata, &e), "cannot unmarshal test data into slice")
	require.Len(t, e, 4, "missing or superflous elements in events slice")

	ctrl := gomock.NewController(t)
	visitor := events.NewMockEventVisitor(ctrl)

	visit := visitor.EXPECT().VisitFailedLogin(MatchJSON(events.FailedLogin{IPAddress: "127.0.0.1"}))
	visit = visitor.EXPECT().VisitLogin(MatchJSON(events.Login{Method: events.LoginUsernamePassword})).After(visit)
	visit = visitor.EXPECT().VisitUpdateAttendance(MatchJSON(events.UpdateAttendance{
		DateID:  123,
		Value:   events.WillNotAttend,
		Comment: "meine Oma hat Geburtstag",
	})).After(visit)
	visitor.EXPECT().VisitLogout(MatchJSON(events.Logout{})).After(visit)

	require.True(t, e.Accept(visitor), "visitor did not match all events")
}

type JSONMatcher struct {
	data []byte
}

// Matches implements gomock.Matcher.
func (j JSONMatcher) Matches(x any) bool {
	data, err := json.Marshal(x)
	if err != nil {
		return false
	}
	return bytes.Equal(data, j.data)
}

// String implements gomock.Matcher.
func (j JSONMatcher) String() string {
	return string(j.data)
}

func MatchJSON(a any) gomock.Matcher {
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return JSONMatcher{data: data}
}
