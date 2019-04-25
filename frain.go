package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	VersionNumber = "[undefined]"

	LogPlain   = 0
	LogDebug   = 1
	LogVerbose = 2
)

// Page specifies the developer tool to check. The Name field here is essentially akin
// to Name field already defined in Service, Component, Incident and IncidentUpdate.
type Page struct {
	Name           string           `json:"name"`
	Services       []Service        `json:"services"`
	Incidents      []Incident       `json:"incidents"`
	IncidentUpdate []IncidentUpdate `json:"incident_update"`
}

// Service represents an external service that we intend to check for availability
type Service struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	PageID        string `json:"page_id"`
	Status        string `json:"status"`
	StatusPageURL string `json:"status_page_url"`
	Provider      string `json:"provider"`

	IsActive bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Components []Component `json:"components"`
}

// Component contains information about a service's components
type Component struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServiceID   string `json:"service_id"`
	ComponentID string `json:"component_id"`
	Status      string `json:"status"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Incident gives all neccessary information relating to a single incident
type Incident struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ServiceID  string `json:"service_id"`
	IncidentID string `json:"incident_id"`
	Status     string `json:"status"`
	Impact     string `json:"impact"`
	Shortlink  string `json:"shortlink"`

	IsActive bool `json:"is_active"`

	ResolvedAt time.Time `json:"resolved_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Events []IncidentUpdate `json:"events"`
}

// IncidentUpdate provides an update to an existing incident
type IncidentUpdate struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	IncidentUpdateID string `json:"incident_update_id"`
	IncidentID       string `json:"incident_id"`
	Status           string `json:"status"`
	Body             string `json:"body"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Init is a simple method to print various build info
func Init() {
	fmt.Println(fmt.Sprintf("Frainserver v%s", VersionNumber))
}

// GetPage returns all information relating to a particular service
func GetPage(pageName string) (*Page, error) {
	return nil, errors.New("this function has not yet been implemented")
}

// GetComponent gets the information relating to a given service component
func GetComponent(compID string) (*Component, error) {
	return nil, errors.New("this function has not yet been implemented")
}

// GetService gets the information relating to a given service
func GetService(serviceID string) (*Service, error) {
	return nil, errors.New("this function has not yet been implemented")
}

// GetIncident gets incident information as well as associated events
func GetIncident(incidentID string) (*Incident, error) {
	return nil, errors.New("this function has not yet been implemented")
}
