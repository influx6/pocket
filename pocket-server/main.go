package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/influx6/faux/web"
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

var context = "pocket-app"

//==============================================================================

func main() {

	app := web.New(events, true, nil, nil)

	go http.ListenAndServe(":3000", app)

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
