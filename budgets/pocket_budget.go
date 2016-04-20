// Package budgets contains the creation code for creating a pocket instance which
// is a fully self enclosed rendering application which can produce its own
// interface markup using the gu library.
package budgets

import (
	"sync/atomic"

	"github.com/influx6/coquery/client"
	"github.com/influx6/gu/gudispatch"
	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
	"github.com/influx6/gu/guviews"
)

//==============================================================================

func init() {
	guviews.Register("pocket/pocket-budget", func(bc BudgetOptions) guviews.Renderable {
		return NewPocketBudget(bc)
	})
}

//==============================================================================

// BudgetOptions defines a configuration struct passed into build initializers.
type BudgetOptions struct {
	UUID     string
	Server   client.Server
	Currency Currency
}

// PocketBudget provides the central repository for creating a pocket instance.
type PocketBudget struct {
	BudgetOptions
	action int64
	active *Budget
	items  map[string]Budget
}

// NewPocketBudget returns a new PocketBudget instance.
func NewPocketBudget(bc BudgetOptions) *PocketBudget {
	pocket := PocketBudget{
		BudgetOptions: bc,
		items:         make(map[string]Budget),
	}

	gudispatch.Subscribe(func(bn *NewBudget) {
		if bc.UUID != bn.UUID {
			return
		}

		// Add new budget into the app list.
		pocket.AddBudget(bn.Title, bn.Price)

		// Dispatch to the view which got registered to update itself.
		gudispatch.Dispatch(guviews.ViewUpdate{ID: bc.UUID})
	})

	return &pocket
}

// AddBudget returns the giving budget with the provided title.
func (p *PocketBudget) AddBudget(title string, budgetPrice float64) *Budget {
	var bu Budget

	atomic.AddInt64(&p.action, 1)
	{
		if bux, ok := p.items[title]; ok {
			bu = bux
		} else {
			bu = Budget{
				Title:    title,
				Price:    budgetPrice,
				currency: p.Currency,
			}

			p.items[title] = bu
		}
	}
	atomic.AddInt64(&p.action, -1)

	return &bu
}

// Render returns the markup defined for a budget.
func (p *PocketBudget) Render() gutrees.Markup {
	var m gutrees.Markup

	atomic.AddInt64(&p.action, 1)
	{

		if p.active != nil {
			m = p.active.Render()
		} else {

			m = elems.Div(attrs.Class("pocket-budget"))

			for _, item := range p.items {
				item.RenderBase().Apply(m)
			}

		}

	}
	atomic.AddInt64(&p.action, -1)

	return m
}
