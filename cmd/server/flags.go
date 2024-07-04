package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	flagEndpointAddr string

	flagOpts FlagOpts
)

type FlagOpts struct {
	SrvAddr netAddress
}

type netAddress struct {
	Host string
	Port int
}

func (n *netAddress) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
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
	n.Host = parts[0]
	n.Port = port

	return nil
}

func parseFlags() {
	flag.StringVar(&flagEndpointAddr, "a", "localhost:8080", "server addr host and port")

	flag.Parse()

	var endpointAddr netAddress
	if err := endpointAddr.Set(flagEndpointAddr); err != nil {
		fmt.Printf("Error parsing endpoint address: %v\n", err)
		os.Exit(1)
	}

	flagOpts = FlagOpts{
		SrvAddr: endpointAddr,
	}
}
