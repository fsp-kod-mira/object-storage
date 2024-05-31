package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Host string `env:"APP_HOST" env-default:"0.0.0.0"`
		Port int    `env:"APP_PORT" env-default:"50051"`
	}

	Logger struct {
		Level string `env:"LOGGER_LEVEL" env-default:"debug"`
	}

	Nats struct {
		Host   string `env:"NATS_HOST"`
		Port   int    `env:"NATS_PORT" env-default:"4222"`
		Bucket string `env:"NATS_BUCKET" env-default:"dev"`
	}
}

func New() *Config {
	config := &Config{}

	if err := cleanenv.ReadEnv(config); err != nil {
		header := "OBJECT STORAGE SERVICE ENVs"
		f := cleanenv.FUsage(os.Stdout, config, &header)
		f()
		panic(err)
	}

	log.Println(config)

	return config
}
