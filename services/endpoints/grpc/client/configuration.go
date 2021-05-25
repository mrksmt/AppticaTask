package client

import "github.com/alexflint/go-arg"

type Config struct {
	Host        string `arg:"env:GRPC_HOST"`
	RequestType string `arg:"env:REQUEST_TYPE"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Host:        "localhost:8080",
		RequestType: "streaming",
	}
	arg.Parse(cfg)
}
