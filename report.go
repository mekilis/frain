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
	Incidents(bool)
	All(bool)
}

// Text is a construct to display the page information in text
type Text struct {
	Data *Page
}

// Incidents implements the Report interface
func (t Text) Incidents(quiet bool) {
	if quiet {
		n := 0
		t1 := time.Now()
		for _, i := range t.Data.Service.Incidents {
			t2 := i.CreatedAt
			if t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year() {
				n++
			}
		}
		fmt.Printf("%d incident(s) reported today.\n", n)
		return
	}

	w := new(tabwriter.Writer)
	printIncidents(w, t.Data.Service.Incidents)
}

// All implements the Report interface
func (t Text) All(quiet bool) {
	w := new(tabwriter.Writer)
	service := t.Data.Service

	name := strings.Title(service.Name)
	titleService := fmt.Sprintf("%s Services", name)
	if quiet {
		summarize(titleService, service.Components, service.Incidents)
		return
	}

	fmt.Printf("----------------------------------\n")
	fmt.Println(titleService)
	printComponents(w, service.Components)
	fmt.Println()
	printIncidents(w, service.Incidents)
}

func summarize(title string, components []Component, incidents []Incident) error {
	op := 0
	numC := 0
	for _, c := range components {
		if c.Status == "operational" {
			op++
		}
		numC++
	}

	fmt.Printf("%s: %d/%d component(s) are operational. %d incident(s) reported.\n",
		title,
		op,
		numC,
		len(incidents),
	)

	return nil
}

func printComponents(w *tabwriter.Writer, comps []Component) {
	colComponents := "\nComponent Name\tStatus\n-------------\t-------"
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, colComponents)
	for _, c := range comps {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", strings.Title(c.Name), strings.Title(c.Status)))
	}
	w.Flush()

	if len(comps) == 0 {
		fmt.Println("No component reports")
	}
}

func printIncidents(w *tabwriter.Writer, inc []Incident) {
	colIncidents := "Date\tTime\tStatus\tDescription\tUpdated\n----------" +
		"\t---------\t----------\t-------------\t------------"

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Println("Incident History\n-----------------")

	fmt.Fprintln(w, colIncidents)

	n := len(inc)
	if n == 0 {
		fmt.Println("No incident reports")
		return
	}

	for _, i := range inc {
		// update using the last incidentUpdate ?? uncertain
		if nIU := len(i.IncidentUpdates); nIU > 0 {
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
				elapsed,
			),
		)
	}
	w.Flush()
}
