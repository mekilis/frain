package frain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// Version represents the current version of this program
	Version string

	// ServiceList stores all supported services on frain
	ServiceList = []string{
		"github",
		"twilio",
		"fastly",
		"bitbucket",
		"circle_ci",
		"status_page",
	}
)

// Page specifies the developer tool to check. The Name field here is essentially akin
// to Name field already defined in Service, Component, Incident and IncidentUpdate.
type Page struct {
	Name    string   `json:"name"`
	Service *Service `json:"service"`
}

// Service represents an external service that we intend to check for availability
type Service struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	PageID        string `json:"pageId"`
	Status        string `json:"status"`
	StatusPageURL string `json:"statusPageUrl"`
	Provider      string `json:"provider"`
	Description   string `json:"description"`
	Indicator     string `json:"indicator"`

	IsActive bool `json:"isActive"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Components []Component `json:"components"`
	Incidents  []Incident  `json:"incidents"`
}

// Component contains information about a service's components
type Component struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServiceID   string `json:"serviceId"`
	ComponentID string `json:"componentId"`
	Status      string `json:"status"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Incident gives all neccessary information relating to a single incident
type Incident struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ServiceID  string `json:"serviceId"`
	IncidentID string `json:"incidentId"`
	Status     string `json:"status"`
	Impact     string `json:"impact"`
	Shortlink  string `json:"shortlink"`

	IsActive bool `json:"isActive"`

	ResolvedAt time.Time `json:"resolvedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	IncidentUpdates []IncidentUpdate `json:"incidentUpdates"`
}

// IncidentUpdate provides an update to an existing incident
type IncidentUpdate struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	IncidentUpdateID string `json:"incidentUpdateId"`
	IncidentID       string `json:"incidentId"`
	Status           string `json:"status"`
	Body             string `json:"body"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Init is a simple method to print various build info
func Init() {
	fmt.Println(fmt.Sprintf("Frain version %s", Version))
	fmt.Println("\nA status checker for various developer tools.")
}

// CleanTimeArg takes in a time construct in string format and ensures it is accurate else
// it returns an error
func CleanTimeArg(t string) (string, error) {
	yMd := strings.Split(t, "-")
	if len(yMd) != 3 {
		return "", errors.New("time must have the format: YYYY-MM-DD")
	}

	y, err := strconv.Atoi(yMd[0])
	if err != nil {
		return "", errors.New(fmt.Sprint("failed to parse year arg in ", t))
	}

	M, err := strconv.Atoi(yMd[1])
	if err != nil {
		return "", errors.New(fmt.Sprint("failed to parse month arg in ", t))
	}

	d, err := strconv.Atoi(yMd[2])
	if err != nil {
		return "", errors.New(fmt.Sprint("failed to parse day arg in ", t))
	}

	return fmt.Sprintf("%04d-%02d-%02d", y, M, d), nil
}
