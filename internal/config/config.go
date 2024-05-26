package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Env         string `yaml:"env" env-default:"local"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"4s"`
	User        string        `yaml:"user" env-required:"true"`
}

func MustLoad() *Config {

	cfgPath := os.Getenv("CFG_PATH")

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("Couldn't load cfg path %v", cfgPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("Couldn't read cfg file %v", err)

	}
	return &cfg

}
