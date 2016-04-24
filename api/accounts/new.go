package accounts

import (
	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/elems"
)

// UserRecord defines a struct which details  a pocket user account.
type UserRecord struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// NewUser defines a struct for creating a new pocket user account.
type NewUser struct{}

// Render returns the markup which renders the new-user account markup.
func (nu *NewUser) Render() gutrees.Markup {
	root := elems.Div()
	return root
}
