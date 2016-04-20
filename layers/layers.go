package layers

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/coquery/client"
	"github.com/influx6/faux/utils"
	"github.com/influx6/gu/gudispatch"
	"github.com/influx6/gu/guviews"
	"github.com/influx6/pocket/budgets"
)

//==============================================================================

// uuis defines a uuid-generator which allows us generate increasing uuid values.
var uuis = utils.NewIncr(0)

//==============================================================================

// PocketLayer returns a view instanced with a Pocket rendering provider.
func PocketLayer(currency string, qs client.Server, mount *js.Object) guviews.Views {

	cu, err := budgets.BudgetCurrency.Find(currency)
	if err != nil {
		gudispatch.Dispatch(&budgets.Notify{
			Message: err.Error(),
			Type:    budgets.BadCurrency,
		})
		return nil
	}

	uuid := uuis.New()

	guviews.MustCreate(guviews.ViewConfig{
		Name:  "pocket/pocket-budget",
		ID:    uuid,
		Paths: []string{"/", "/pockets"},
		Param: budgets.BudgetOptions{
			UUID:     uuid,
			Server:   qs,
			Currency: cu,
		},
	})

	view := guviews.MustGet(uuid)
	view.Mount(mount)

	return view
}

//==============================================================================

// LoginLayer instantiates the login layer for the application, setting up
// and returning the view concerned with login.
func LoginLayer(qs client.Server, mount *js.Object) guviews.Views {

	return nil
}

//==============================================================================

// AccountLayer instantiates the account layer for the application, setting up
// and returning the view concerned with login.
func AccountLayer(qs client.Server, mount *js.Object) guviews.Views {

	return nil
}
