package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/dawsonalex/aggregator/lib"

	"github.com/google/uuid"

	"github.com/dawsonalex/aggregator/watcher"
	"github.com/sirupsen/logrus"
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

		node := struct {
			Instance uuid.UUID
			port     int
		}{}
		err = json.Unmarshal(body, &node)
		if err != nil {
			log.Errorf("Error unmarshalling node: %v", err)
			http.Error(w, "Error parding JSON", http.StatusBadRequest)
			return
		}
		reg.AddNode(node.Instance)
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
		log.WithField("ID", nodeID).Println("Removing node")
		reg.RemoveNode(nodeID)
	})
}

// FilesHandler handles requests to the /files endpoint.
func FilesHandler(reg *watcher.Registry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// sort and return the filelist
			files := reg.ListFiles()
			sort.Strings(files)
			fileResponse := struct {
				files []string
			}{
				files: files,
			}
			log.Println("Responding with file list.")
			json.NewEncoder(w).Encode(fileResponse)

		} else if r.Method == http.MethodPatch {
			// Do the file operation
			defer r.Body.Close()
			decoder := json.NewDecoder(r.Body)
			var operations lib.OperationRequests

			err := decoder.Decode(&operations)
			if err != nil {
				log.Errorf("Error decoding JSON: %v", err)
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}

			for _, op := range operations {
				log.WithFields(logrus.Fields{
					"ID": op.Instance,
					"op": op.Type,
				}).Println("Doing file operation")
				if node := reg.Node(op.Instance); node != nil {
					node.Do(watcher.Operation{
						Type:     op.Type,
						SeqNo:    op.SeqNo,
						Filename: op.Value.Filename,
					})
				}
			}
		}
	})
}
