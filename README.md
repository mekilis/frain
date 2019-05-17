# Frain CLI

Frain makes it possible to check the availability of many developer tools or services in the terminal.

Services currently supported are:
* GitHub
* Fastly
* Twilio

## Installation

### Building from source
You can get the entire source code via the `go` tool using the following:

`$ go get github.com/mekilis/frain/cmd/...`

Move into project directory via:

`$ cd $GOPATH/src/github.com/mekilis/frain`

Installation could then be done in the `$GOPATH/bin` folder using the command:

`$ make install`

OR

The compiled binary could be moved to another location via:

`$ make build && sudo mv frain /usr/bin` 

on Linux, for example.

**NOTE:** Building from source requires Go (version 1.11 or later). It is assumed that `$GOPATH` is set.

### Precompiled binaries
You can download precompiled binaries from the [release page](https://github.com/mekilis/frain/releases).

You can download them easily with the following:

```bash
sudo curl -L -o /usr/local/bin/frain https://github.com/mekilis/frain/releases/download/v0.1.0/frain-v0.1.0-<os-type>-amd64 && sudo chmod +x /usr/local/bin/frain

# Windows
sudo curl -L -o /usr/local/bin/frain https://github.com/mekilis/frain/releases/download/v0.1.0/frain-v0.1.0-windows-amd64.exe && sudo chmod +x /usr/local/bin/frain
```

The `os-type` can be either linux or darwin (for Mac)

## Usage
```
Usage:
        frain [options] <args>...

        -c <path>,      --config=<path>         Specifies path to configuration file with a
                                                list of services to check
        -f <format>,    --format=<format>       Specifies result output format i.e. txt, json
                                                or xml (txt by default)
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
        frain github incidents 2019-01-12 2019-05-05    ==> Fetch incidents from start to end dates
```

