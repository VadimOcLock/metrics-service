package config

type Agent struct {
	AppConfig
	AgentConfig
}

type WebServer struct {
	AppConfig
	WebServerConfig
	BackupConfig
	DatabaseConfig
}
