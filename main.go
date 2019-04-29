package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	format  = "Select format to display query"
	help    = "Displays this help"
	version = "Current version of frain"

	formatFlag  = flag.String("format", "txt", format)
	helpFlag    = flag.Bool("help", false, help)
	versionFlag = flag.Bool("version", false, version)
)

func init() {
	flag.StringVar(formatFlag, "f", "txt", format)
	flag.BoolVar(helpFlag, "h", false, help)
	flag.BoolVar(versionFlag, "v", false, version)

	flag.Usage = func() {
		usage := "Usage:\n\tfrain [options] service\n\nOptions:\n" +
			"\t-c, --config=CONFIGPATH\tSpecifies path to configuration file\n" +
			"\t-f, --format=FORMAT\tSpecifies result output format i.e. txt, json or xml (txt by default)\n" +
			"\t-h, --help\t\tDisplays this message\n" +
			"\t-q, --quiet\t\tDisplays just the summary\n" +
			"\t-v, --version\t\tDisplays the current version of program\n"
		fmt.Print(usage)
	}
}

func main() {
	flag.Parse()

	if *versionFlag {
		Init()
		os.Exit(0)
	}

	if *helpFlag {
		Init()
		flag.Usage()
		os.Exit(0)
	}

	argsLen := len(os.Args)
	if argsLen < 2 {
		fmt.Println("Error: no service specified.")
		flag.Usage()
		os.Exit(1)
	}

	format := strings.ToLower(*formatFlag)
	if format != "txt" && format != "json" && format != "xml" {
		fmt.Printf("Error: bad format specified '%s'\n", format)
		flag.Usage()
		os.Exit(1)
	}

	var page Page
	name := os.Args[1]
	page.Name = name

	services, err := Services()
	if err != nil {
		fmt.Println("Error: failed to get page names")
		os.Exit(2)
	}

	if _, ok := services[name]; !ok {
		fmt.Printf("Error: unknown service specified '%s'\n", name)
		os.Exit(2)
	}
	page.Services = services[name] // TODO: check if deep copy

	incidents, err := Incidents(name)
	if err != nil {
		fmt.Printf("Error: failed to get incidents for %s\n", name)
		os.Exit(2)
	}
	page.Incidents = make([]Incident, 0)
	if _, ok := incidents[name]; ok {
		page.Incidents = incidents[name]
	}

	var report Report
	switch format {
	case "json":
		//pass
	case "xml":
		//pass
	default:
		report = Text{
			Data: &page,
		}
	}

	err = report.All()
	if err != nil {
		fmt.Println("Error: failed to generate report -", err)
		os.Exit(3)
	}
}
