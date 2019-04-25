package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Report
type Report interface {
	Incident(*Incident) error
	Service(*Service) error
	All() error
}

type Text struct {
	Data *Page
}

func (t Text) Incident(i *Incident) error {
	return nil
}

func (t Text) Service(s *Service) error {
	return nil
}

func (t Text) All() error {
	w := new(tabwriter.Writer)

	name := strings.Title(t.Data.Name)
	titleServices := fmt.Sprintf("%s Services", name)
	colServices := "\nName\tStatus\n-----\t-------"

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Println(titleServices)
	fmt.Fprintln(w, colServices)
	for _, s := range t.Data.Services {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", s.Name, s.Status))
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

	return nil
}
