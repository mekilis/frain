package frain

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

var (
	pageDoesNotExist = errors.New("specified service name does not exist")
)

type Data struct {
	All []Service `json:"getAllServices"`
}

type Result struct {
	Data `json:"data"`
}

func Services() (map[string][]Service, error) {
	query := bytes.NewBuffer([]byte(os.Getenv("QUERY_ALL_SERVICES")))
	host := os.Getenv("FRAIN_HOST")
	jsn := "application/json"

	resp, err := http.Post(host, jsn, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	queryMap := make(map[string][]Service)
	for _, service := range result.Data.All {
		queryMap[service.Name] = append(queryMap[service.Name], service)
	}

	return queryMap, nil
}

func Incidents(name string) (map[string][]Incident, error) {
	return make(map[string][]Incident), nil
}

func IncidentUpdates(name string, i Incident) (map[string][]IncidentUpdate, error) {
	return make(map[string][]IncidentUpdate), nil
}
