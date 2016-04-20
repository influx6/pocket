package budgets

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
)

// Budget defines a collection of cost items writting against a given pocket
// budget, it hosts the central items for a budget.
type Budget struct {
	Title            string  `json:"title"`
	Price            float64 `json:"price"`
	items            []BudgetItem
	currency         Currency
	action           int64
	activeBudgetItem int
}

// AddItem adds a new budget item into the lists of Budgets.
func (b *Budget) AddItem(title string, desc string, price float64) {
	bi := BudgetItem{
		Title:  title,
		Desc:   desc,
		Price:  price,
		Time:   time.Now(),
		Budget: b,
	}

	atomic.AddInt64(&b.action, 1)
	{
		b.items = append(b.items, bi)
	}
	atomic.AddInt64(&b.action, -1)
}

// RenderBase returns a markup to render the basic view of a Budget.
func (b *Budget) RenderBase() gutrees.Markup {
	root := elems.Div(
		attrs.Class("budget"),
		elems.Label(attrs.Class("budget-item-price"), elems.Text(fmt.Sprintf("%s%.2f", b.currency, b.Price))),
		elems.Label(attrs.Class("budget-item-count"), elems.Text(fmt.Sprintf("%d", len(b.items)))),
	)
	return root
}

// Render returns the markup defined for a budget item which is to be rendered.
func (b *Budget) Render() gutrees.Markup {
	root := elems.Div(attrs.Class("budget"))

	barView := elems.Div(attrs.Class("budget-bar", "side-left"))
	barItems := elems.Div(attrs.Class("budget-items", "side-right"))

	for _, item := range b.items {
		item.Render().Apply(barItems)
	}

	barView.Apply(root)
	barItems.Apply(root)

	return root
}
