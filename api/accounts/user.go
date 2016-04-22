package accounts

import (
	"fmt"

	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
	"github.com/influx6/pocket/api/currency"
)

// UserCurrency defines the pockets associated with a specific currency.
type UserCurrency struct {
	Currency currency.Currency
	pocketID string
}

// Render returns the rendereable markup for the UserCurrency struct.
func (u *UserCurrency) Render() gutrees.Markup {
	root := elems.Div(attrs.Class("user-currency"))
	elems.Label(attrs.Class("user-currency-name"), elems.Text(u.Currency.Name)).Apply(root)
	elems.Label(attrs.Class("user-currency-sign"), elems.Text(u.Currency.Name)).Apply(root)
	return root
}

// User defines the logged in user who owns the current records.
type User struct {
	accounts []UserCurrency
}

// Render returns the rendereable markup for the User struct.
func (u *User) Render() gutrees.Markup {
	root := elems.Div(attrs.Class("user", "pocket-user", "pocket-user-account"))

	for ind, uc := range u.accounts {
		rn := uc.Render()
		attrs.ID(fmt.Sprintf("currency-item-#%d", ind)).Apply(rn)
		rn.Apply(root)
	}

	return root
}
