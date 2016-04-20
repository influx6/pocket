package auth

import (
	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
)

// Login provides the login authentication provider for log-in into
// a budget system.
type Login struct{}

// Render returns the markup for the login page.
func (l Login) Render() gutrees.Markup {
	return elems.Div(
		attrs.Class("login", "full-page"),
	)
}
