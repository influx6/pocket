package budgets

import (
	"fmt"
	"time"

	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
)

// BudgetItem defines a price item which defines a subcost to a given Budget in
// a pocket.
type BudgetItem struct {
	Title  string    `json:"title"`
	Desc   string    `json:"desc"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
	Budget *Budget   `json:"budget"`
}

// Render returns the markup defined for a BudgetItem which is to be rendered.
func (b *BudgetItem) Render() gutrees.Markup {
	var tag string

	if dlen := len(b.Title); dlen > 3 {
		tag = b.Title[:3]
	} else {
		tag = b.Title
	}

	return elems.Div(
		attrs.Class("budget-item"),
		elems.Label(attrs.Class("budget-item-price"), elems.Text(fmt.Sprintf("%s%.2f", b.Budget.currency, b.Price))),
		elems.Label(attrs.Class("budget-item-name"), elems.Text(fmt.Sprintf("%s..", tag))),
	)
}
