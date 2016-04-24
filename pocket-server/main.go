package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
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

var contexts = "pocket-app"

//==============================================================================

func main() {

	pocketapp := app.New(events, true, nil, nil)

	app.PageRoute(pocketapp, "GET", "/", func(ctx context.Context, w *app.ResponseRequest, params app.Param) error {

		return nil
	})

	go http.ListenAndServe(":3000", pocketapp)

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
