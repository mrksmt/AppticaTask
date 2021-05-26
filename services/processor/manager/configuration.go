package manager

import (
	"log"
	"task/api"

	"github.com/alexflint/go-arg"
)

type Config struct {
	api.CommonParams
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

	if cfg.ApplicationId <= 0 {
		log.Fatalf("Wrong Application Id: %d", cfg.ApplicationId)
	}
	if cfg.CountryId <= 0 {
		log.Fatalf("Wrong Country Id: %d", cfg.CountryId)
	}
	if cfg.UpdateRate <= 1 {
		log.Fatalf("Wrong Update Rate: %d", cfg.UpdateRate)
	} else {
		log.Printf("Manager update rate: %d sec", cfg.UpdateRate)
	}

}
