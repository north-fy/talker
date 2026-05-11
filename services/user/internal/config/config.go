package config

type PostgresCfg struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
}
