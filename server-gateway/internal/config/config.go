package config

import (
	"go-simpler.org/env"
)

type Config struct {
	HTTPCfg       HTTPCfg       `env:"HTTP_"`
	UserSrvCfg    UserSrvCfg    `env:"USER_SRV_"`
	MessageSrvCfg MessageSrvCfg `env:"MESSAGE_SRV_"`
	ChatSrvCfg    ChatSrvCfg    `env:"CHAT_SRV_"`
	JWTSecret     string        `env:"JWT_SECRET"`
}

type HTTPCfg struct {
	Port int `env:"PORT"`
}

type UserSrvCfg struct {
	Addr string `env:"ADDR"`
}

type MessageSrvCfg struct {
	Addr string `env:"ADDR"`
}

type ChatSrvCfg struct {
	Addr string `env:"ADDR"`
}

func (c *Config) Load() error {
	return env.Load(c, nil)
}
