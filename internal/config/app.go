package config

type AppConfig struct {
	// Health, Log, Prometheus, etc.
}

type AgentConfig struct {
	EndpointAddr   string `env:"ADDRESS"`
	PoolInterval   int    `env:"REPORT_INTERVAL"`
	ReportInterval int    `env:"POLL_INTERVAL"`
}

type WebServerConfig struct {
	SrvAddr string `env:"ADDRESS"`
}
