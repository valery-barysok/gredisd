package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/gredisd-app"
)

var usageStr = `
Usage: gredisd [options]

Server Options:
    -a, --addr <host>                Bind to host address (default: 0.0.0.0)
    -p, --port <port>                Use port for clients (default: 16379)
        --databases <count>          Set the number of databases. The default database is DB 0, you can select
                                     a different one on a per-connection basis using SELECT <dbid> where
                                     dbid is a number between 0 and 'databases'-1
        --trace_protocol             Trace low level read/write operations

Authorization Options:
        --auth <token>               Authorization token required for connections

Common Options:
    -h, --help                       Show this message
    -v, --version                    Show version
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func main() {
	opts := app.Options{}
	initFromEnv(&opts)

	var showVersion bool

	flag.IntVar(&opts.Port, "port", opts.Port, "Port to listen on.")
	flag.IntVar(&opts.Port, "p", opts.Port, "Port to listen on.")
	flag.StringVar(&opts.Host, "addr", opts.Host, "Network host to listen on.")
	flag.StringVar(&opts.Host, "a", opts.Host, "Network host to listen on.")
	flag.BoolVar(&opts.TraceProtocol, "trace_protocol", false, "Trace low level read/write operations")
	flag.StringVar(&opts.Auth, "auth", "", "Password for AUTH command.")
	flag.IntVar(&opts.Databases, "databases", app.DefaultDatabases, "Password for AUTH command.")
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.BoolVar(&showVersion, "v", false, "Print version information.")

	flag.Usage = usage

	flag.Parse()

	gApp := gredisd.NewApp(&opts)
	if showVersion {
		gApp.ShowVersion()
		os.Exit(0)
	}

	gApp.Run()
}

func initFromEnv(opts *app.Options) {
	intEnvVar(&opts.Port, "GREDIS_PORT", app.DefaultPort)
	stringVar(&opts.Host, "GREDIS_HOST", app.DefaultHost)

}

func intEnvVar(p *int, name string, value int) {
	val := os.Getenv(name)
	if val != "" {
		v, err := strconv.Atoi(val)
		if err == nil {
			*p = v
			return
		}
	}
	*p = value
}

func stringVar(p *string, name string, value string) {
	val := os.Getenv(name)
	if val != "" {
		*p = val
		return
	}
	*p = value
}
