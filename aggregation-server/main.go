package main

import (
	"net/http"

	"github.com/dawsonalex/aggregator/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting aggregator")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", server.HelloHandler())
	mux.HandleFunc("/bye", server.ByeHandler())
	mux.HandleFunc("files", server.FilesHandler())
}
