package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const HTTPProtocolName = "http"

var (
	flagEndpointAddr   string
	flagReportInterval int
	flagPoolInterval   int

	flagOpts FlagOpts
)

type FlagOpts struct {
	EndpointAddr   netAddress
	PoolInterval   time.Duration
	ReportInterval time.Duration
}

type netAddress struct {
	Protocol string
	Host     string
	Port     int
}

func (n *netAddress) String() string {
	return fmt.Sprintf("%s://%s:%d", n.Protocol, n.Host, n.Port)
}

func (n *netAddress) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address format, expected host:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}
	n.Protocol = HTTPProtocolName
	n.Host = parts[0]
	n.Port = port

	return nil
}

func parseFlags() {
	flag.IntVar(&flagReportInterval, "r", 10, "report frequency in seconds")
	flag.IntVar(&flagPoolInterval, "p", 2, "poll data frequency in seconds")
	flag.StringVar(&flagEndpointAddr, "a", "localhost:8080", "server endpoint host and port")

	flag.Parse()

	var endpointAddr netAddress
	if err := endpointAddr.Set(flagEndpointAddr); err != nil {
		fmt.Printf("Error parsing endpoint address: %v\n", err)
		os.Exit(1)
	}

	flagOpts = FlagOpts{
		EndpointAddr:   endpointAddr,
		PoolInterval:   time.Duration(flagPoolInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
	}
}
