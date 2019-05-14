# Frain CLI
CLI for checking the availability of various developer tools. 

## Installation
`go get github.com/mekilis/frain/cmd/...`

## Usage
```bash
Usage:
        frain [options] <args>...

        -c <path>,      --config=<path>         Specifies path to configuration file with a list of services to check
        -f <format>,    --format=<format>       Specifies result output format i.e. txt, json or xml (txt by default)
        -h,             --help                  Displays this help message
        -q <service>,   --quiet <service>       Displays just the summary for specified service
        -v,             --version               Displays the current version of this program

Args:
        <service>
        <service> incidents
        <service> incidents <start time> <end time>

Note that both start and end times have the format YYYY-MM-DD

Examples:
        frain github                                    ==> Fetch report for github
        frain -q github                                 ==> Summarize fetched result for github
        frain github incidents                          ==> Fetch only incident reports
        frain github incidents 2019-01-12               ==> Fetch incidents from start date
        frain github incidents 2019-01-12 2019-05-05    ==> Fetch incidents from start to end date
```
