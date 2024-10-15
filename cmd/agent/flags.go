package main

import (
	"flag"
	"fmt"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"os"
	"strconv"
	"strings"

	"github.com/VadimOcLock/metrics-service/internal/config"
)

const (
	defaultSrvAddr        = "localhost:8080"
	defaultPoolInterval   = 2
	defaultReportInterval = 10
	HTTPProtocolName      = "http"
)

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
		return errorz.ErrInvalidAddressFormat
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	n.Host = parts[0]
	n.Port = port

	return nil
}

func parseFlags(cfg *config.Agent) error {
	var (
		flagEndpointAddr   string
		flagReportInterval int
		flagPoolInterval   int
	)

	flag.IntVar(&flagReportInterval, "r", defaultReportInterval, "report frequency in seconds")
	flag.IntVar(&flagPoolInterval, "p", defaultPoolInterval, "poll data frequency in seconds")
	flag.StringVar(&flagEndpointAddr, "a", defaultSrvAddr, "server endpoint host and port")

	flag.Parse()

	var endpointAddr netAddress
	if err := endpointAddr.Set(flagEndpointAddr); err != nil {
		return fmt.Errorf("error parsing endpoint address: %w", err)
	}

	if envVal := os.Getenv("ADDRESS"); envVal == "" {
		cfg.AgentConfig.EndpointAddr = endpointAddr.String()
	}
	if envVal := os.Getenv("REPORT_INTERVAL"); envVal == "" {
		cfg.AgentConfig.ReportInterval = flagReportInterval
	}
	if envVal := os.Getenv("POLL_INTERVAL"); envVal == "" {
		cfg.AgentConfig.PoolInterval = flagPoolInterval
	}

	return nil
}
