package frain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Data struct {
	All []Service `json:"getAllServices"`
}

type Result struct {
	Data `json:"data"`
}

func Services() (map[string][]Service, error) {
	query := bytes.NewBuffer([]byte("{ \"query\": \"{getAllServices {id, name, statusPageUrl," +
		"provider, indicator, isActive, createdAt, updatedAt, components{id, name, status, " +
		"description}}}\"}"))
	// TODO: readd 'description' in query
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
