package config

type AppConfig struct {
	// Health, Log, Prometheus, etc.
}

type AgentConfig struct {
	EndpointAddr   string `env:"ADDRESS"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

type WebServerConfig struct {
	SrvAddr string `env:"ADDRESS"`
	FileWriter
}

type FileWriter struct {
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}
