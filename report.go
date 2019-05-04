package frain

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Report is an interface implemented by types that generates report in different formats.
type Report interface {
	Incident(*Incident) error
	Service(*Service) error
	All() error
}

// Text is a construct to display the page information in text
type Text struct {
	Data *Page
}

// Incident implements the Report interface
func (t Text) Incident(i *Incident) error {
	return nil
}

// Service implements the Report interface
func (t Text) Service(s *Service) error {
	return nil
}

// All implements the Report interface
func (t Text) All() error {
	w := new(tabwriter.Writer)
	service := t.Data.Service
	fmt.Printf("----------------------------------\n")

	name := strings.Title(service.Name)
	titleService := fmt.Sprintf("%s Services", name)
	colComponents := "\nService Name\tStatus\n-------------\t-------"

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Println(titleService)
	fmt.Fprintln(w, colComponents)
	for _, comp := range service.Components {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", strings.Title(comp.Name), strings.Title(comp.Status)))
	}
	w.Flush()
	if len(service.Components) == 0 {
		fmt.Println("No service reports")
	}

	colIncidents := "Date\tTime\tStatus\tDescription\tUpdated\n----------" +
		"\t---------\t----------\t-------------\t------------"
	n := len(service.Incidents)

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Println("\nIncident History\n-----------------")

	if n > 0 {
		fmt.Fprintln(w, colIncidents)
		for _, i := range service.Incidents {
			// update using the last incidentUpdate ?? uncertain
			if n := len(i.IncidentUpdates); n > 0 {
				i.Status = i.IncidentUpdates[0].Status
				i.UpdatedAt = i.IncidentUpdates[0].UpdatedAt
			}

			elapsed, _ := TimeAgo(i.UpdatedAt, time.Now())
			if elapsed == "0 seconds ago" {
				elapsed = "     -"
			}

			fmt.Fprintln(
				w,
				fmt.Sprintf("%s %d %d\t%d:%d:%d\t%s\t%s\t%s",
					i.CreatedAt.Month(),
					i.CreatedAt.Day(),
					i.CreatedAt.Year(),

					i.CreatedAt.Hour(),
					i.CreatedAt.Minute(),
					i.CreatedAt.Second(),

					i.Status,
					i.Impact,
					elapsed, //TODO: call TimeAgo() here
				),
			)
		}
	} else {
		fmt.Println("No incident reports")
	}
	w.Flush()

	return nil
}
