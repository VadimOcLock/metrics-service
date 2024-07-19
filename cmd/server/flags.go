package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/VadimOcLock/metrics-service/internal/config"
)

const defaultSrvAddr = "localhost:8080"

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
		return errors.New("invalid address format, expected host:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	n.Host = parts[0]
	n.Port = port

	return nil
}

func parseFlags(cfg *config.WebServer) error {
	var flagSrvAddr string

	flag.StringVar(&flagSrvAddr, "a", defaultSrvAddr, "server addr host and port")

	flag.Parse()

	var srvAddr netAddress
	if err := srvAddr.Set(flagSrvAddr); err != nil {
		return fmt.Errorf("error parsing server address: %w", err)
	}

	if envVal := os.Getenv("ADDRESS"); envVal == "" {
		cfg.WebServerConfig.SrvAddr = srvAddr.String()
	}

	return nil
}
