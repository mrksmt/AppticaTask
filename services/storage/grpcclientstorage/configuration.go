package grpcclientstorage

import (
	"log"

	"github.com/alexflint/go-arg"
)

type Config struct {
	DataHost string `arg:"env:DATA_HOST"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		DataHost: "localhost:8095",
	}
	arg.Parse(cfg)
	log.Printf("Dataservice host: %s", cfg.DataHost)
}
