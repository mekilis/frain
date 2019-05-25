package frain

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	color "gopkg.in/gookit/color.v1"
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

var (
	yellow = color.FgLightYellow.Render

	bg = color.BgGreen.Render
	fg = color.FgBlack.Render
)

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

	color.Bold.Println(titleService)
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
	colComponents := "\nCOMPONENT NAME\tSTATUS"
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	fmt.Fprint(w, bg(fg(colComponents)))
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

func printIncidents(w *tabwriter.Writer, inc []Incident) {
	colIncidents := "\nDATE\tTIME\tIMPACT\tUPDATED\tSTATUS"

	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	color.Bold.Println("Incident History")

	n := len(inc)
	if n == 0 {
		fmt.Println("No incident reports")
		return
	}

	fmt.Fprint(w, bg(fg(colIncidents)))

	for j := n - 1; j >= 0; j-- {
		i := inc[j]
		elapsed, _ := TimeAgo(i.UpdatedAt, time.Now())
		if elapsed == "0 seconds ago" {
			elapsed = "     -"
		}

		fmt.Fprint(
			w,
			fmt.Sprintf("\n%s %d %d\t%d:%d:%d\t%s\t%s\t%s",
				i.CreatedAt.Month(),
				i.CreatedAt.Day(),
				i.CreatedAt.Year(),

				i.CreatedAt.Hour(),
				i.CreatedAt.Minute(),
				i.CreatedAt.Second(),

				strings.Title(i.Impact),
				elapsed,
				render(strings.Title(i.Status)),
			),
		)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func render(status string) string {
	var r = color.FgWhite.Render // default
	s := strings.Split(strings.ToLower(status), " ")
	switch s[0] {
	case "operational", "resolved":
		r = color.FgGreen.Render
	case "degraded", "under", "under_maintenance", "investigating":
		r = color.FgYellow.Render
	case "outage", "critical":
		r = color.FgRed.Render
	}

	return r(status)
}
