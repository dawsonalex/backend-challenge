package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dawsonalex/aggregator/lib"

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

		// unmarshal the request into the expected struct.
		var node lib.HelloRequest
		err = json.Unmarshal(body, &node)
		if err != nil {
			log.Errorf("Error unmarshalling node: %v", err)
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		// Add the node and set it's initial files if it's new.
		if fileChan, _, isNew := reg.AddNode(node.Instance); isNew {
			files := make([]string, 0)
			url, err := alterAddress(r.RemoteAddr, node.Port)
			if err == nil {
				// Request the nodes file list
				files, err = watcher.GetNodeFiles(url)
				if err != nil {
					log.Errorf("error getting files from watcher: %v", err)
					return
				}

				// send the files down the file channel
				// for the node to add.
				for _, v := range files {
					fileChan <- v
				}
				close(fileChan)
			} else {
				log.Errorf("error parsing url: %v", err)
				return
			}
		}
	})
}

// Take a remote address, format it, set the port, and return a *url.URL
// that represents the formatted address.
func alterAddress(remoteAddr string, port int) (*url.URL, error) {
	host := strings.Split(remoteAddr, ":")[0]
	newAddr := fmt.Sprintf("http://%s:%d/files", host, port)
	url, err := url.Parse(newAddr)
	if err != nil {
		return nil, err
	}
	return url, nil
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

		var nodeInstance lib.ByeRequest
		err := decoder.Decode(&nodeInstance)
		if err != nil {
			log.Error("Error decoding JSON")
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		nodeID, err := uuid.Parse(nodeInstance.Instance)
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
			fileResponse := lib.FilesResponse{
				Files: files,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
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
				log.WithFields(log.Fields{
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

// Sort filenames in ascending order based on name without extension.
// Uses natural sort.
func sortFiles(filenames []string) []string {
	sort.SliceStable(filenames, func(i, j int) bool {
		// trim the file extension before comparing file names.
		first := strings.TrimSuffix(filenames[i], filepath.Ext(filenames[i]))
		second := strings.TrimSuffix(filenames[j], filepath.Ext(filenames[j]))
		return first < second
	})
	return filenames
}
