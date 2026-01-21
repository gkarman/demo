package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(path string) *Config {
	var cfg Config

	// 1️⃣ Загружаем файл (если есть)
	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
	}

	// 2️⃣ ENV всегда поверх
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read env: %v", err)
	}

	return &cfg
}
