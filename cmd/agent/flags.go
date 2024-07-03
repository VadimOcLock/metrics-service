package main

import (
	"flag"
)

var (
	flagRunAddr        string
	flagReportInterval int
	flagPoolInterval   int
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report frequency in second")
	flag.IntVar(&flagPoolInterval, "p", 2, "pool data frequency in second")

	flag.Parse()
}
