package server

import "github.com/alexflint/go-arg"

type Config struct {
	Port string `arg:"env:GRPC_SRV_PORT"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Port: "8090",
	}
	arg.MustParse(cfg)
}
