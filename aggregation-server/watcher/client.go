package watcher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GetNodeFiles makes a request to a watcher node for
// its file list.
func GetNodeFiles(url *url.URL) ([]string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	filesBody, err := ioutil.ReadAll(resp.Body)

	fileResponse := struct {
		Files []struct {
			Filename string
		}
	}{}
	err = json.Unmarshal(filesBody, &fileResponse)
	if err != nil {
		return nil, errors.New("error reading node response")
	}

	files := make([]string, 0)

	for _, file := range fileResponse.Files {
		files = append(files, file.Filename)
	}
	return files, nil
}
