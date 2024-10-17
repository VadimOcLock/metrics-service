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
	Restore         bool   `env:"RESTORE"`
	Interval        int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_DSN"`
}

func (m *DatabaseConfig) InMemoryMode() bool {
	return m.DSN == ""
}
