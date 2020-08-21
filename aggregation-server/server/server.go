package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/dawsonalex/aggregator/watcher"
	log "github.com/sirupsen/logrus"
)

// HelloHandler handles requests to the /hello endpoint.
func HelloHandler(reg *watcher.Registry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodPost) {
			log.Errorf("Invalid HTTP method, got: %v", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Error reading body: %v", err)
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}

		var node watcher.Node
		err = json.Unmarshal(body, &node)
		if err != nil {
			log.Errorf("Error unmarshalling node: %v", err)
			http.Error(w, "Error parding JSON", http.StatusBadRequest)
			return
		}

		reg.AddNode(&node)
	})
}

// ByeHandler handles requests to the /bye endpoint.
func ByeHandler(reg *watcher.Registry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodPost) {
			log.Errorf("Invalid HTTP method, got: %v", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		nodeInstance := struct {
			ID string `json:"instance"`
		}{}
		err := decoder.Decode(&nodeInstance)
		if err != nil {
			log.Error("Error decoding JSON")
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		nodeID, err := uuid.Parse(nodeInstance.ID)
		if err != nil {
			log.WithField("value", nodeID).Error("Error parsing node ID")
			http.Error(w, "Error parsing node ID", http.StatusBadRequest)
			return
		}
		reg.RemoveNode(nodeID)
	})
}

// FilesHandler handles requests to the /files endpoint.
func FilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
