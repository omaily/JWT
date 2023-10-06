package config

import (
	"log/slog"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Port        string        `yaml:"port" env-default:":4000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

var cfg Config
var once sync.Once

func MustLoad() *Config {
	once.Do(func() {
		slog.Info("read application configuration")
		cfg = Config{}
		if err := cleanenv.ReadConfig("config/conf.yaml", &cfg); err != nil {
			slog.Error("cannot read config: %s", err)
		}
	})
	return &cfg
}
