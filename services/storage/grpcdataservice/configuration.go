package dataservice

import "github.com/alexflint/go-arg"

type Config struct {
	Port string `arg:"env:DATA_SRV_PORT"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Port: "8095",
	}
	arg.Parse(cfg)
}
