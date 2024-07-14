package events

import (
	"github.com/jpicht/polyjson"
)

type Login struct {
	polyjson.Implements[Event]

	Method LoginMethod `json:"method"`
}

type FailedLogin struct {
	polyjson.Implements[Event]

	IPAddress string `json:"ip_address"`
}

type Logout struct {
	polyjson.Implements[Event]
}

type LoginMethod string

const (
	LoginUsernamePassword = LoginMethod("password")
	LoginCookie           = LoginMethod("cookie")
	LoginAdminImpersonate = LoginMethod("impersonate")
)
