package frain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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

var (
	errHTTPPost   = errors.New("Error: a network error occurred while fetching data")
	errJSONDecode = errors.New("Error: failed to decode fetched data")
)

// GetService sends a POST request to the host server and then returns all information
// relating to a developer tool to check
func GetService(name string, startTime, endTime time.Time) (*Service, error) {
	q := fmt.Sprintf(`{"query": "{getService(name:%s)`+
		`{id, name, statusPageUrl, provider, indicator, isActive, createdAt, updatedAt,`+
		` components`+
		`{id, name, status, description},`+
		` incidents(startTime:\"%s\", endTime:\"%s\")`+
		`{id, name,impact, status, isActive, createdAt, shortlink, updatedAt, resolvedAt, incidentUpdates{id, body, status, createdAt, updatedAt}},`+
		` highLevelComponents`+
		`{id, name, status, description}}}"}`, name, parseDate(&startTime), parseDate(&endTime))
	query := bytes.NewBuffer([]byte(q))
	host := os.Getenv("FRAIN_HOST")
	if host == "" {
		host = "https://frain-server.herokuapp.com/graphql"
	}

	resp, err := http.Post(host, "application/json", query)
	if err != nil {
		return nil, errHTTPPost
	}
	defer resp.Body.Close()

	var result SingleResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errJSONDecode
	}

	return &result.Service, nil
}

func parseDate(t *time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
}

// GetServiceList returns a list of services currently supported by frain
func GetServiceList() ([]string, error) {
	query := bytes.NewBuffer([]byte(`{ "query": "{getAllServices {name}}" }`))
	host := os.Getenv("FRAIN_HOST")
	if host == "" {
		host = "https://frain-server.herokuapp.com/graphql"
	}

	resp, err := http.Post(host, "application/json", query)
	if err != nil {
		return nil, errHTTPPost
	}
	defer resp.Body.Close()

	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errJSONDecode
	}

	var services []string
	var sMap = map[string]bool{}

	for _, s := range result.All {
		sMap[strings.ToLower(s.Name)] = true
	}

	for s := range sMap {
		services = append(services, s)
	}

	return services, nil
}
