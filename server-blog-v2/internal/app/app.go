package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"server-blog-v2/config"
)

// Run 运行应用。
func Run(cfg *config.Config) {
	app, cleanup, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("app - Run - InitializeApp: %v", err)
	}
	defer cleanup()

	app.HTTPServer.Start()
	app.Logger.Info("app - Run - started: %s v%s", app.Info.Name, app.Info.Version)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		app.Logger.Info("app - Run - signal: %s", s.String())
	case err = <-app.HTTPServer.Notify():
		app.Logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = app.HTTPServer.Shutdown()
	if err != nil {
		app.Logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	app.Logger.Info("app - Run - stopped: %s v%s", app.Info.Name, app.Info.Version)
}
