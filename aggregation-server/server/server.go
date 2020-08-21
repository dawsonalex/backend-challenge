package server

import (
	"net/http"
)

// HelloHandler handles requests to the /hello endpoint.
func HelloHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

// ByeHandler handles requests to the /bye endpoint.
func ByeHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

// FilesHandler handles requests to the /files endpoint.
func FilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
