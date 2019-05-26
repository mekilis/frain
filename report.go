package frain

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	wordwrap "github.com/mitchellh/go-wordwrap"

	"github.com/fatih/color"
)

// Report is an interface implemented by types that generates report in different formats.
type Report interface {
	Incidents(bool, bool)
	All(bool, bool)
}

// Text is a construct to display the page information in text
type Text struct {
	Data *Page
}

var (
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgHiGreen)

	titleBar = color.New(color.FgBlack, color.BgWhite)
	bold     = color.New(color.Bold)
)

const maxWidth = 40

// Incidents implements the Report interface
func (t Text) Incidents(quiet, full bool) {
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
	printIncidents(w, t.Data.Service.Incidents, full)
}

// All implements the Report interface
func (t Text) All(quiet, full bool) {
	w := new(tabwriter.Writer)
	service := t.Data.Service

	sb := strings.Builder{}
	words := strings.Split(service.Name, "_")
	if len(words) > 0 {
		words[0] = strings.Title(words[0])
	}
	for _, w := range words {
		sb.Write([]byte(w))
		sb.Write([]byte(" "))
	}
	name := strings.TrimSpace(sb.String())

	titleService := fmt.Sprintf("%s Services", name)
	if quiet {
		summarize(titleService, service.Components, service.Incidents)
		return
	}

	bold.Println(titleService)
	printComponents(w, service.Components)
	fmt.Println()
	printIncidents(w, service.Incidents, full)
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
	colComponents := "\nCOMPONENT NAME\tSTATUS"
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	titleBar.Fprint(w, colComponents)
	for _, c := range comps {
		words := strings.Split(c.Status, "_")
		sb := strings.Builder{}
		for _, word := range words {
			sb.WriteString(strings.Title(word))
			sb.WriteString(" ")
		}
		status := strings.TrimSpace(sb.String())

		fmt.Fprint(w, fmt.Sprintf("\n%s\t%s", strings.Title(c.Name), render(status)))
	}
	fmt.Fprintln(w)
	w.Flush()

	if len(comps) == 0 {
		fmt.Println("No component reports")
	}
}

func printIncidents(w *tabwriter.Writer, inc []Incident, full bool) {
	colIncidents := "\nDATE\tTIME\tIMPACT\tUPDATED\tDESCRIPTION\tSTATUS\t"

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	bold.Println("Incident History")

	n := len(inc)
	if n == 0 {
		fmt.Println("No incident reports")
		return
	}

	titleBar.Fprint(w, colIncidents)

	for j := n - 1; j >= 0; j-- {
		i := inc[j]
		elapsed, _ := TimeAgo(i.UpdatedAt, time.Now())
		if elapsed == "0 seconds ago" {
			elapsed = "     -"
		}

		description := "-"
		for _, x := range i.IncidentUpdates {
			if i.Status == x.Status {
				description = x.Body
				break
			}
		}

		dte := fmt.Sprintf("%s %d, %d",
			i.CreatedAt.Month(),
			i.CreatedAt.Day(),
			i.CreatedAt.Year(),
		)

		tme := fmt.Sprintf("%d:%d:%d",
			i.CreatedAt.Hour(),
			i.CreatedAt.Minute(),
			i.CreatedAt.Second(),
		)

		desc := wrap(description, maxWidth)
		n := len(desc)
		if !full {
			others := false
			if n > 1 {
				others = true
			}
			desc[0] = pad(desc[0], maxWidth, others)
		}

		fmt.Fprint(
			w,
			fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t",
				dte,
				tme,

				strings.Title(i.Impact),
				elapsed,
				desc[0],
				render(strings.Title(i.Status)),
			),
		)
		if full {
			n := len(desc)
			for i := 1; i < n; i++ {
				fmt.Fprint(w, "\n\t\t\t\t", desc[i], "\t\t")
			}
		}
	}
	fmt.Fprintln(w)
	w.Flush()
}

func render(status string) string {
	var r = color.New()
	s := strings.Split(strings.ToLower(status), " ")
	switch s[0] {
	case "operational", "resolved":
		r.Add(color.FgGreen)
	case "degraded", "under", "under_maintenance", "investigating":
		r.Add(color.FgYellow)
	case "outage", "critical":
		r.Add(color.FgRed)
	default:
		r.Add(color.FgWhite)
	}

	return r.Sprint(status)
}

func wrap(s string, width int) []string {
	if width < 0 {
		return []string{"-"}
	}

	s = wordwrap.WrapString(s, uint(width))
	words := strings.Split(s, "\n")
	n := len(words)
	if n == 0 {
		return []string{"-"}
	}

	return words
}

// This pads any string s with three dots ('.') for a given pad length
func pad(s string, padLength int, others bool) string {
	if n := len(s); padLength <= n || (!others && n < padLength) {
		return s
	}

	s += strings.Repeat(".", padLength)
	s = s[:padLength]
	dots := 0
	n := len(s)

	for j := n - 1; j >= 0; j-- {
		if s[j] != '.' {
			break
		}
		dots++
	}

	if dots < 3 {
		s = s[:padLength-3] + "..."
	} else {
		s = s[:padLength-dots+3]
	}

	return s
}
