package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/dawsonalex/aggregator/watcher"

	"github.com/dawsonalex/aggregator/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting aggregator")

	reg := watcher.NewRegistry()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", server.HelloHandler(reg))
	mux.HandleFunc("/bye", server.ByeHandler(reg))
	mux.HandleFunc("/files", server.FilesHandler(reg))

	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait here until SIGINT received, then exec callback function
	// to gracefully shutdown.
	awaitInterrupt(func(done chan bool) {
		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
		done <- true
	})
	log.Info("Aggregator stopped.")
}

func awaitInterrupt(onInterrupt func(chan bool)) {
	done := make(chan bool)
	go func() {
		// Wait for SIGINT to stop services.
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		defer signal.Stop(sigchan)
		<-sigchan

		go onInterrupt(done)
	}()

	<-done
}
