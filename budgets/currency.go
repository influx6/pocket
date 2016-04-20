package budgets

import (
	"fmt"
	"strings"
)

//==============================================================================

// Currency defines the currency type for a budget.
type Currency struct {
	Name string `json:"name"`
	Sign string `json:"sign"`
}

// String returns the sign for the given currency.
func (c Currency) String() string {
	return c.Sign
}

//==============================================================================

// Currencies declares a lists of currency slice.
type Currencies []Currency

// Find returns a currency by its giving name else returns an error if
// not found.
func (c Currencies) Find(cm string) (Currency, error) {
	var fc Currency
	var found bool

	for _, cu := range c {
		if strings.ToLower(cu.Name) == strings.ToLower(cm) {
			fc = cu
			found = true
			break
		}
	}

	if !found {
		return fc, fmt.Errorf("Unknown Currency[%s]", cm)
	}

	return fc, nil
}

// BudgetCurrency defines a lists of currency types and their respective signs.
var BudgetCurrency = Currencies{
	Currency{
		Name: "Dollars",
		Sign: "$",
	},
	Currency{
		Name: "Naira",
		Sign: "#",
	},
}

//==============================================================================
