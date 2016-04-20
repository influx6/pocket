package budgets

//==============================================================================

// Alert defines the types of alerts defined within the pocket app.
type Alert int

const (
	// InvalidUUID defines an alert for when an invalid uuid is provided for a
	// pocket app.
	InvalidUUID Alert = iota + 1

	// BadCurrency is defined for when a unknown/invalid currency is provided for
	// use.
	BadCurrency
)

//==============================================================================

// Notify defines the struct for sending notifications for budget actions.
type Notify struct {
	Message string
	Type    Alert
}

//==============================================================================

// SyncBudget is used to send a sync signal that a budget should resync
// its data.
type SyncBudget struct {
	UUID string
}

// NewBudget defines a struct for requesting the creation of a new budget item.
type NewBudget struct {
	By    string
	UUID  string
	Title string
	Price float64
}

// AmendBudget defines a struct for requesting the amendation of a budget item.
type AmendBudget struct {
	By    string
	UUID  string
	Title string
	Price float64
}

//==============================================================================
