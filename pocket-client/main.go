package main

import (
	"fmt"
	"time"

	"github.com/influx6/coquery/client"
	"github.com/influx6/coquery/client/js"
	"github.com/influx6/pocket/layers"

	"honnef.co/go/js/dom"
)

//==============================================================================

var events eventlog

// logg provides a concrete implementation of a logger.
type eventlog struct{}

// Log logs all standard log reports.
func (l eventlog) Log(context interface{}, name string, message string, data ...interface{}) {
	fmt.Printf("Log: %s : %s : %s : %s\n", context, "DEV", name, fmt.Sprintf(message, data...))
}

// Error logs all error reports.
func (l eventlog) Error(context interface{}, name string, err error, message string, data ...interface{}) {
	fmt.Printf("Error: %s : %s : %s : %s : Error %s\n", context, "DEV", name, fmt.Sprintf(message, data...), err)
}

//==============================================================================

func main() {
	window := dom.GetWindow()
	doc := window.Document()

	client := client.NewServo(events, "http://127.0.0.1:3000", 300*time.Millisecond, js.HTTP)

	layers.AccountLayer(client, doc.QuerySelector("body").Underlying())

}
