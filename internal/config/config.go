package config

import "fmt"

type Config struct {
	Env    string `env:"ENV" env-default:"local"`
	DB     DBConfig
	Logger LoggerConfig
	ServerHttp ServerHttpConfig
}

type DBConfig struct {
	Host                         string `env:"DB_HOST"`
	Port                         int    `env:"DB_PORT"`
	User                         string `env:"DB_USER"`
	Password                     string `env:"DB_PASS"`
	Name                         string `env:"DB_NAME"`
	SSLMode                      string `env:"DB_SSLMODE"`
	MaxConnections               int32  `env:"DB_MAX_CONNECTIONS"`
	MinConnections               int32  `env:"DB_MIN_CONNECTIONS"`
	MaxConnectionLifeTimeMinutes int    `env:"DB_MAX_CONNECTION_LIFETIME_MINUTES"`
	MaxConnectionIdleTimeMinutes int    `env:"DB_MAX_CONNECTION_IDLE_TIME_MINUTES"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

type ServerHttpConfig struct {
	Addr                string `env:"SERVER_HTTP_ADDR" env-default:":8080"`
	ReadTimeoutSeconds  int    `env:"SERVER_HTTP_READ_TIMEOUT_SECONDS" env-default:"10"`
	WriteTimeoutSeconds int    `env:"SERVER_HTTP_WRITE_TIMEOUT_SECONDS" env-default:"10"`
}
