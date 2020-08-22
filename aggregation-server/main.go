package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/dawsonalex/aggregator/watcher"

	"github.com/dawsonalex/aggregator/server"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort     = 8000
	defaultLogLevel = "info"
)

func main() {
	var logLevel = flag.String("log", defaultLogLevel, "the level of logging (debug, info, warning, fatal, panic)")
	var port = flag.Uint("p", defaultPort, "the port to listen on")
	flag.Parse()

	log := initLogger(*logLevel)
	log.Info("Starting aggregator")

	reg := watcher.NewRegistry(log)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", server.HelloHandler(reg))
	mux.HandleFunc("/bye", server.ByeHandler(reg))
	mux.HandleFunc("/files", server.FilesHandler(reg))

	addr := fmt.Sprintf(":%d", *port)
	log.Info("listening on port: ", *port)
	srv := &http.Server{
		Addr:    addr,
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
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
		done <- true
	})
	log.Info("Aggregator stopped.")
}

func initLogger(logLevel string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.Out = os.Stdout

	if logLevel, err := logrus.ParseLevel(logLevel); err != nil {
		fmt.Printf("error during log init: %v\n", err)
		fmt.Println("using default log level 'info'")
		log.Level = logrus.InfoLevel
	} else {
		fmt.Println("using log level: ", logLevel.String())
		log.Level = logLevel
	}
	return log
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
