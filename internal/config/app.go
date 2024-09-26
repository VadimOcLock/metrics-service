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
}

type BackupConfig struct {
	Restore         bool   `env:"RESTORE" envDefault:"true"`
	Interval        int    `env:"STORE_INTERVAL" envDefault:"300"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./backups/storage.txt"`
}
