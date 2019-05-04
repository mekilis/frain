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

type SingleData struct {
	Service Service `json:"getService"`
}

type SingleResult struct {
	SingleData `json:"data"`
}

func Services() (map[string][]Service, error) {
	query := bytes.NewBuffer([]byte("{ \"query\": \"{getAllServices {id, name, statusPageUrl," +
		"provider, indicator, isActive, createdAt, updatedAt, components{id, name, status, " +
		"description}, incidents{id, name,impact, status, isActive, createdAt, shortlink, updatedAt," +
		"resolvedAt, incidentUpdates{id, body, status, createdAt, updatedAt}}}}\"}"))
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

func GetServiceFor(name string) (*Service, error) {
	query := bytes.NewBuffer([]byte("{ \"query\": \"{getService(name:" + name +
		") {id, name, statusPageUrl," +
		"provider, indicator, isActive, createdAt, updatedAt, components{id, name, status, " +
		"description}, incidents(startTime:\\\"2019-05-01\\\", endTime:\\\"2019-05-05\\\"){id, name,impact, status, isActive, createdAt, shortlink, updatedAt," +
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
