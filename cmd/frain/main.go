package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/mekilis/frain"
)

var (
	green  = color.New(color.FgGreen).Sprint
	yellow = color.New(color.FgYellow).Sprint

	config  = "Path to configuration file"
	format  = "Select format to display query"
	help    = "Displays this help"
	full    = "Displays a full version of incident descriptions"
	list    = "Lists the currently supported services"
	quiet   = "Displays the service summary"
	version = "Current version of frain"

	configFlag  = flag.String("config", "", config)
	formatFlag  = flag.String("format", "txt", format)
	helpFlag    = flag.Bool("help", false, help)
	fullFlag    = flag.Bool("full", false, full)
	listFlag    = flag.Bool("list", false, list)
	quietFlag   = flag.Bool("quiet", false, quiet)
	versionFlag = flag.Bool("version", false, version)

	buildVersion string

	subCommands = map[string]bool{
		"incidents": true,
	}
)

func init() {
	flag.StringVar(configFlag, "c", "", config)
	flag.StringVar(formatFlag, "f", "txt", format)
	flag.BoolVar(helpFlag, "h", false, help)
	flag.BoolVar(listFlag, "l", false, list)
	flag.BoolVar(quietFlag, "q", false, quiet)
	flag.BoolVar(versionFlag, "v", false, version)

	flag.Usage = func() {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 2, 8, 0, '\t', tabwriter.AlignRight)
		fmt.Fprint(w, yellow("\nUsage:"),
			"\n\tfrain ", green("[options]"), " <args>...\n",
			yellow("\nOptions:"),
			green("\n\t\t--full\t"), "Displays the full version of incident descriptions",
			green("\n\t-c <path>,\t--config=<path>\t"), "Specifies path to configuration file with a\n\t\t\tlist of services to check",
			green("\n\t-f <format>,\t--format=<format>\t"), "Specifies result output format i.e. txt, json\n\t\t\tor xml (txt by default)",
			green("\n\t-h,\t--help\t"), "Displays this help message",
			green("\n\t-l,\t--list\t"), "Lists the currently supported services on frain",
			green("\n\t-q <service>,\t--quiet <service>\t"), "Displays just the summary for specified service",
			green("\n\t-v,\t--version\t"), "Displays the current version of this program\n",
			yellow("\nArgs:"),
			"\n\t<service>\n\t<service> ", green("incidents"),
			"\n\t<service> ", green("incidents <start time> <end time>\n\n"),
			"Note that both start and end times have the format YYYY-MM-DD\n",
			yellow("\nExamples:"),
			"\n\tfrain github\t==> Fetch report for github",
			"\n\tfrain -q github\t==> Summarize fetched result for github",
			"\n\tfrain github incidents\t==> Fetch only incident reports",
			"\n\tfrain github incidents 2019-01-12\t==> Fetch incidents from start date",
			"\n\tfrain github incidents 2019-01-12 2019-05-05\t==> Fetch incidents from start to end dates\n")

		w.Flush()
	}
}

func main() {
	flag.Parse()
	setupVerInfo()
	parseFlagOptions()

	flagArgs := flag.Args()
	format := strings.ToLower(*formatFlag)

	if format != "txt" && format != "json" && format != "xml" {
		fmt.Printf("frain: bad format specified '%s' (\"frain help\" for help)\n", format)
		exit()
	}

	if len(os.Args) < 2 || len(flagArgs) == 0 {
		fmt.Println("frain: no service specified (\"frain help\" for help)")
		exit()
	}

	subCommand, startTime, endTime, err := parseFlagParams(flagArgs)
	if err != nil {
		fmt.Println("frain:", err, "(\"frain help\" for help)")
		exit()
	}

	name := strings.ToLower(flagArgs[0])
	report, err := getReport(name, format, startTime, endTime)
	if err != nil {
		fmt.Println("frain:", err)
		exit()
	}

	if len(flagArgs) < 2 {
		report.All(*quietFlag, *fullFlag)
	} else {

		switch subCommand {

		case "incidents":
			report.Incidents(*quietFlag, *fullFlag)

		}
	}
}

func setupVerInfo() {
	// go build -ldflags "-X main.buildVersion=X.Y.Z"
	if buildVersion == "" {
		buildVersion = "<undefined>"
	}
	frain.Version = buildVersion
}

func parseFlagOptions() {
	if *versionFlag {
		frain.Init()
		exit()
	}

	if *helpFlag {
		frain.Init()
		flag.Usage()
		exit()
	}

	if *listFlag {
		frain.Init()
		listServices()
		exit()
	}

	if len(*configFlag) != 0 {
		fmt.Println("-c flag set. this feature has not yet been implemented. please check again later.")
		exit()
	}
}

func progress(c chan int) {
	s := "Please wait while fetching data"
	dots := []string{".  ", ".. ", "..."}
	for {
		select {
		case <-c:
			return
		default:
			for _, d := range dots {
				fmt.Print(s, d)
				time.Sleep(time.Second)
				fmt.Print("\r \r")
			}

		}
	}
}

func clear() {
	cls := "                                        "
	fmt.Print("\r \r")
	fmt.Print(cls)
	fmt.Print("\r \r")
}

func listServices() {
	var c = make(chan int)
	go progress(c)

	sl, err := frain.GetServiceList()
	c <- 0
	clear()
	if err != nil {
		fmt.Println(err)
		return
	}

	size := len(sl)
	if size == 0 {
		fmt.Println("\nNo service supported at the moment.")
		return
	}

	fmt.Println("\nServices currently supported are:")
	for _, s := range sl {
		fmt.Printf("\t%s\n", s)
	}
}

func parseFlagParams(flagArgs []string) (string, time.Time, time.Time, error) {
	subCommand := ""
	startTime, _ := time.Parse("2006-01-02", "1970-01-01") // iso layout
	endTime := time.Now()
	flagArgsLen := len(flagArgs)

	if flagArgsLen > 1 {
		subCommand = flagArgs[1]
		if _, ok := subCommands[subCommand]; !ok {
			return subCommand, startTime, endTime, fmt.Errorf("unknown subcommand specified for %s: '%s'", flagArgs[0], flagArgs[1])
		}
	}

	if flagArgsLen > 2 {
		sTime, err := frain.CleanTimeArg(flagArgs[2])
		if err != nil {
			return subCommand, startTime, endTime, fmt.Errorf("start time error. %v", err)
		}
		startTime, err = time.Parse("2006-01-02", sTime)
		if err != nil {
			return subCommand, startTime, endTime, fmt.Errorf("bad format specified for start time: %v", flagArgs[2])
		}
	}

	if flagArgsLen > 3 {
		eTime, err := frain.CleanTimeArg(flagArgs[3])
		if err != nil {
			return subCommand, startTime, endTime, fmt.Errorf("end time error. %v", err)
		}
		endTime, err = time.Parse("2006-01-02", eTime)
		if err != nil {
			return subCommand, startTime, endTime, fmt.Errorf("bad format specified for end time: %v", flagArgs[3])
		}
	}

	return subCommand, startTime, endTime, nil
}

func getReport(name, format string, startTime, endTime time.Time) (frain.Report, error) {
	var report frain.Report
	var page frain.Page
	errFmt := "report feature has not been implemented yet, please check back later in a future release"

	switch format {
	case "json":
		return nil, fmt.Errorf("json %v", errFmt)

	case "xml":
		return nil, fmt.Errorf("xml %v", errFmt)

	case "txt":
		report = frain.Text{
			Data: &page,
		}
	}

	var c = make(chan int)
	go progress(c)

	service, err := frain.GetService(name, startTime, endTime)
	c <- 1
	clear()
	if err != nil {
		return nil, err
	}

	if strings.ToLower(service.Name) != name {
		return nil, fmt.Errorf("'%s' is not a recognized service on frain (see \"frain --list\")", name)
	}

	page.Name, page.Service = name, service

	return report, nil
}

func exit() {
	// other cleanup tasks
	os.Exit(0)
}
