package main

import (
	"flag"
	"fmt"
	"os"
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
}
