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

const (
	defaultSrvAddr         = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "./store_metrics.txt"
	defaultRestore         = false
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
	var (
		flagSrvAddr         string
		flagStoreInterval   int
		flagFileStoragePath string
		flagRestore         bool
	)

	flag.StringVar(&flagSrvAddr, "a", defaultSrvAddr, "server addr host and port")
	flag.IntVar(&flagStoreInterval, "i", defaultStoreInterval, "interval store save to file")
	flag.StringVar(&flagFileStoragePath, "f", defaultFileStoragePath, "path to store save file")
	flag.BoolVar(&flagRestore, "r", defaultRestore, "restore metrics in file")

	flag.Parse()

	var srvAddr netAddress
	if err := srvAddr.Set(flagSrvAddr); err != nil {
		return fmt.Errorf("error parsing server address: %w", err)
	}

	if envVal := os.Getenv("ADDRESS"); envVal == "" {
		cfg.WebServerConfig.SrvAddr = srvAddr.String()
	}
	if envVal := os.Getenv("STORE_INTERVAL"); envVal == "" {
		cfg.WebServerConfig.FileWriter.StoreInterval = flagStoreInterval
	}
	if envVal := os.Getenv("FILE_STORAGE_PATH"); envVal == "" {
		cfg.WebServerConfig.FileWriter.FileStoragePath = flagFileStoragePath
	}
	if envVal := os.Getenv("RESTORE"); envVal == "" {
		cfg.WebServerConfig.FileWriter.Restore = flagRestore
	}

	return nil
}
