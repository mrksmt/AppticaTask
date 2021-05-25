package manager

import (
	"log"

	"github.com/alexflint/go-arg"
)

type Config struct {
	UpdateRate int `arg:"env:UPDATE_RATE"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		UpdateRate: 10,
	}
	arg.Parse(cfg)
	log.Printf("Manager update rate: %d sec", cfg.UpdateRate)
}
