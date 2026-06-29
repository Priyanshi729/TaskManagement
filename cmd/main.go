package main

import (
	"Task-Management/database"
	"Task-Management/server"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const shutdownTimOut = 10 * time.Second

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SetupRoutes()

	if err := database.ConnectDB(); err != nil {
		log.Panicf("failed to connect to DB with err: %v", err)
	}
	logrus.Print("migration successful!!")

	go func() {
		if err := srv.Run(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :8080")

	<-done

	logrus.Info("shutting down server")

	if err := srv.Shutdown(shutdownTimOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}

	if err := database.CloseDB(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

}
