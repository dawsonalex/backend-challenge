package main

import (
	"net/http"

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
	mux.HandleFunc("files", server.FilesHandler())
}
