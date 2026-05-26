package config

import "go-simpler.org/env"

type Config struct {
	PostgresCfg PostgresCfg `env:"POSTGRES_"`
	GRPCCfg     GRPCCfg     `env:"GRPC_"`
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

func (c *Config) Load() error {
	return env.Load(c, nil)
}
