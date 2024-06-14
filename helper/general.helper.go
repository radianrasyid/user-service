package helper

import (
	"fmt"
	"net/http"
)

func GetCurrentConnectionAndLoad(url string) (int, float64, error) {
	// You can replace this URL with the URL of your service endpoint

	// Send HTTP GET request to the service
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0.0, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return 0, 0.0, fmt.Errorf("failed to get connection and load info. Status code: %d", resp.StatusCode)
	}

	// Parse the response body to get connection and load information
	// Here, we're assuming the response is in JSON format
	var data struct {
		Connections int     `json:"connections"`
		Load        float64 `json:"load"`
	}

	return data.Connections, data.Load, nil
}
