package manager

import "github.com/alexflint/go-arg"

type Config struct {
	Rate int `arg:"env:UPDATE_RATE"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Rate: 10,
	}
	arg.MustParse(cfg)
}
