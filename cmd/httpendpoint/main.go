package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task/application"
	"task/services/endpoints/http/enpoints/apptopcat"
	"task/services/endpoints/http/server"
	"task/services/storage/grpcclientstorage"

	"github.com/pkg/errors"
)

func main() {

	log.SetFlags(log.Lshortfile)

	// make and build the app
	appName := "AppTopCategory HTTP endpoint"
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

	// var storage api.Storage
	// storage = new(mockStorage.ErrorStorage)
	// storage = new(mockStorage.EmptyStorage)
	// storage = new(mockStorage.AnyDateStorage)
	// storage := ledisstorage.New()

	stor := grpcclientstorage.New(nil)
	app.Add(stor)

	appTopCategoryEndpoint := apptopcat.New(stor)

	httpServer := server.New(nil, appTopCategoryEndpoint)
	app.Add(httpServer)
}
