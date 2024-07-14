package typefield

import "github.com/jpicht/polyjson"

type Login struct {
	polyjson.Implements[Event] `polytypeid:"login"`

	Method LoginMethod `json:"method"`
}

type FailedLogin struct {
	polyjson.Implements[Event] `polytypeid:"login_failed"`

	IPAddress string `json:"ip_address"`
}

type Logout struct {
	polyjson.Implements[Event] `polytypeid:"logout"`
}

type LoginMethod string

const (
	LoginUsernamePassword = LoginMethod("password")
	LoginCookie           = LoginMethod("cookie")
	LoginAdminImpersonate = LoginMethod("impersonate")
)
