package frain

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
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

	// range through all services
	for _, service := range t.Data.Services {

		name := strings.Title(service.Name)
		titleService := fmt.Sprintf("%s Services", name)
		colComponents := "\nName\tStatus\n-----\t-------"

		w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
		fmt.Println(titleService)
		fmt.Fprintln(w, colComponents)
		for _, comp := range service.Components {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s", strings.Title(comp.Name), strings.Title(comp.Status)))
		}
		w.Flush()

		colIncidents := "\nDate\tTime\tStatus\tDescription\tUpdated\n"
		n := len(t.Data.Incidents)

		w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
		fmt.Println("\nIncident History")

		if n > 0 {
			fmt.Fprintln(w, colIncidents)
			for _, i := range t.Data.Incidents {
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
						i.UpdatedAt, //TODO: call TimeAgo() here
					),
				)
			}
		} else {
			fmt.Println("No incidents reports")
		}
		w.Flush()

		fmt.Printf("----------------------------------\n\n") // for next service
	}

	return nil
}
