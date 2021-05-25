package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task/application"
	"task/services/processor/manager"
	"task/services/processor/processor"
	"task/services/processor/source"
	"task/services/storage/grpcclientstorage"

	"github.com/pkg/errors"
)

func main() {

	log.SetFlags(log.Lshortfile)

	// make and build the app
	appName := "AppTopCategory Data Processor"
	appVersion := "v0.1"
	app := application.New(appName, appVersion)
	buildApp(app)

	// run the app
	if err := app.Run(); err != nil {
		log.Fatal(errors.Wrap(err, "app.Run"))
	}

	// waiting for signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// stop and waiting for shutdown
	app.Stop()
}

// make all app modules and add some to app Run
func buildApp(app *application.Application) {

	src := source.New()
	proc := processor.New()

	// stor := ledisstorage.New()
	// app.Add(stor)

	stor := grpcclientstorage.New(nil)
	app.Add(stor)

	mgr := manager.New(src, proc, stor, nil)
	app.Add(mgr)
}
