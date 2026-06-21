package config

import (
	"time"

	"go-simpler.org/env"
)

type Config struct {
	PostgresCfg PostgresCfg `env:"POSTGRES_"`
	GRPCCfg     GRPCCfg     `env:"GRPC_"`
	RedisCfg    RedisCfg    `env:"REDIS_"`
}

type PostgresCfg struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"DB"`
}

type GRPCCfg struct {
	Port int `env:"PORT"`
}

type RedisCfg struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	ReadTimeout time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
	DB   int `env:"DB"`
}

func (c *Config) Load() error {
	return env.Load(c, nil)
}
