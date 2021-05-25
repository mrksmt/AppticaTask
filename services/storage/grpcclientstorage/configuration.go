package grpcclientstorage

import "github.com/alexflint/go-arg"

type Config struct {
	Host string `arg:"env:DATA_SRV_HOST"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Host: "localhost:8095",
	}
	arg.MustParse(cfg)
}
