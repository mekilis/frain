package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mekilis/frain"
)

var (
	config  = "Path to configuration file"
	format  = "Select format to display query"
	help    = "Displays this help"
	list    = "Lists the currently supported services"
	quiet   = "Displays the service summary"
	version = "Current version of frain"

	configFlag  = flag.String("config", "", config)
	formatFlag  = flag.String("format", "txt", format)
	helpFlag    = flag.Bool("help", false, help)
	listFlag    = flag.Bool("list", false, list)
	quietFlag   = flag.Bool("quiet", false, quiet)
	versionFlag = flag.Bool("version", false, version)

	buildVersion string
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
		w.Init(os.Stdout, 4, 8, 0, '\t', tabwriter.AlignRight)
		fmt.Fprint(w, "\nUsage:\n\tfrain [options] <args>...\n\n")
		fmt.Fprintln(w, "\t-c <path>,\t--config=<path>\tSpecifies path to configuration file with a ")
		fmt.Fprintln(w, "\t\t\tlist of services to check")
		fmt.Fprintln(w, "\t-f <format>,\t--format=<format>\tSpecifies result output format i.e. txt, json")
		fmt.Fprintln(w, "\t\t\tor xml (txt by default)")
		fmt.Fprintln(w, "\t-h,\t--help\tDisplays this help message")
		fmt.Fprintln(w, "\t-l,\t--list\tLists the currently supported services on frain")
		fmt.Fprintln(w, "\t-q <service>,\t--quiet <service>\tDisplays just the summary for specified service")
		fmt.Fprintln(w, "\t-v,\t--version\tDisplays the current version of this program")
		fmt.Fprintln(w, "\nArgs:")
		fmt.Fprintln(w, "\t<service>\n\t<service> incidents")
		fmt.Fprint(w, "\t<service> incidents <start time> <end time>\n\n")
		fmt.Fprintln(w, "Note that both start and end times have the format YYYY-MM-DD")
		fmt.Fprintln(w, "\nExamples:")
		fmt.Fprintln(w, "\tfrain github\t==> Fetch report for github")
		fmt.Fprintln(w, "\tfrain -q github\t==> Summarize fetched result for github")
		fmt.Fprintln(w, "\tfrain github incidents\t==> Fetch only incident reports")
		fmt.Fprintln(w, "\tfrain github incidents 2019-01-12\t==> Fetch incidents from start date")
		fmt.Fprintln(w, "\tfrain github incidents 2019-01-12 2019-05-05\t==> Fetch incidents from start to end dates")

		w.Flush()
	}
}

func main() {
	flag.Parse()
	versionInfo()

	if *versionFlag {
		frain.Init()
		os.Exit(0)
	}

	if *helpFlag {
		frain.Init()
		flag.Usage()
		os.Exit(0)
	}

	if *listFlag {
		frain.Init()
		listServices()
		os.Exit(0)
	}

	if len(*configFlag) != 0 {
		fmt.Println("-c flag set. this feature has not yet been implemented. please check again later.")
		os.Exit(0)
	}

	format := strings.ToLower(*formatFlag)
	if format != "txt" && format != "json" && format != "xml" {
		fmt.Printf("Error: bad format specified '%s'\n", format)
		flag.Usage()
		os.Exit(1)
	}

	osArgsLen := len(os.Args)
	flagArgs := flag.Args()
	flagArgsLen := len(flagArgs)
	if osArgsLen < 2 || flagArgsLen == 0 {
		fmt.Println("Error: no service specified.")
		flag.Usage()
		os.Exit(1)
	}

	startTime, _ := time.Parse("2006-01-02", "1970-01-01") // iso layout
	endTime := time.Now()

	var err error
	if flagArgsLen > 1 {
		if flagArgs[1] != "incidents" {
			fmt.Printf("unknown query specified for %s: '%s'\n", flagArgs[0], flagArgs[1])
			os.Exit(2)
		}
		if flagArgsLen > 2 {
			sTime, err := frain.CleanTimeArg(flagArgs[2])
			if err != nil {
				fmt.Println("start time error.", err)
				os.Exit(4)
			}
			startTime, err = time.Parse("2006-01-02", sTime)
			if err != nil {
				fmt.Println("bad format specified for start time:", sTime)
				os.Exit(4)
			}
		}
		if flagArgsLen > 3 {
			eTime, err := frain.CleanTimeArg(flagArgs[3])
			if err != nil {
				fmt.Println("end time error.", err)
				os.Exit(4)
			}
			endTime, err = time.Parse("2006-01-02", eTime)
			if err != nil {
				fmt.Println("bad format specified for end time:", eTime)
				os.Exit(4)
			}
		}

		// other subcommands
	}

	var page frain.Page
	name := strings.ToLower(flagArgs[0])
	page.Name = name

	var c = make(chan int)
	go progress(c)

	service, err := frain.GetService(name, startTime, endTime)
	c <- 1
	clear()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if strings.ToLower(service.Name) != name {
		fmt.Printf("Error: '%s' is not a recognized service on frain. Try running 'frain --list'.\n", name)
		os.Exit(2)
	}

	page.Service = service

	var report frain.Report
	errFmt := "report feature has not been implemented yet, please check back later in a future release."
	switch format {
	case "json":
		fmt.Println("json", errFmt)
		os.Exit(0)
	case "xml":
		fmt.Println("xml", errFmt)
		os.Exit(0)
	case "txt":
		report = frain.Text{
			Data: &page,
		}
	default:
		fmt.Println("bad file format specified:", format)
		os.Exit(1)
	}

	if flagArgsLen < 2 {
		report.All(*quietFlag)
	} else {
		report.Incidents(*quietFlag)
	}
}

func versionInfo() {
	// go build -ldflags "-X main.buildVersion=X.Y.Z"
	if buildVersion == "" {
		buildVersion = "<undefined>"
	}
	frain.Version = buildVersion
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
