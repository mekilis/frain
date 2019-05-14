package frain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Data is a model to match the `getAllServices` JSON tag from frain backend
type Data struct {
	All []Service `json:"getAllServices"`
}

// Result is a model to match the `data` JSON tag from frain backend
type Result struct {
	Data `json:"data"`
}

// SingleData is a model to match the `getService` JSON tag from frain backend for a single service
type SingleData struct {
	Service Service `json:"getService"`
}

// SingleResult is a model to match the `data` JSON tag from frain backend for a single service
type SingleResult struct {
	SingleData `json:"data"`
}

// GetService sends a POST request to the host server and then returns all information
// relating to a developer tool to check
func GetService(name string, startTime, endTime time.Time) (*Service, error) {
	query := bytes.NewBuffer([]byte("{ \"query\": \"{getService(name:" + name + ")" +
		"{id, name, statusPageUrl, provider, indicator, isActive, createdAt, updatedAt, components{id, name, status, description}," +
		"incidents(startTime:\\\"" + parseDate(&startTime) + "\\\", endTime:\\\"" +
		parseDate(&endTime) + "\\\"){id, name,impact, status, isActive, createdAt, shortlink, updatedAt," +
		"incidentUpdates{id, body, status, createdAt, updatedAt}}}}\"}"))
	// TODO: fix resolvedAt issues
	host := os.Getenv("FRAIN_HOST")
	if host == "" {
		host = "https://frain-server.herokuapp.com/graphql"
	}
	jsn := "application/json"

	resp, err := http.Post(host, jsn, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SingleResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Service, nil
}

func parseDate(t *time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
}
