package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mekilis/frain"
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
		frain.Init()
		os.Exit(0)
	}

	if *helpFlag {
		frain.Init()
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

	var page frain.Page
	name := strings.ToLower(os.Args[1])
	page.Name = name

	service, err := frain.GetServiceFor(name)
	if err != nil {
		fmt.Println("Error: failed to get service information for", name)
		log.Println(err)
		os.Exit(2)
	}

	if service.Name != name {
		fmt.Printf("Error: unknown service specified '%s'\n", name)
		os.Exit(2)
	}

	page.Service = service

	var report frain.Report
	switch format {
	case "json":
		//pass
	case "xml":
		//pass
	default:
		report = frain.Text{
			Data: &page,
		}
	}

	err = report.All()
	if err != nil {
		fmt.Println("Error: failed to generate report -", err)
		os.Exit(3)
	}
}
