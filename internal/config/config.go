package config

type Config struct {
	Env string `env:"ENV" env-default:"local"`
	DB  DBConfig
}

type DBConfig struct {
	Host string `env:"DB_HOST"`
	Port int    `env:"DB_PORT"`
	Name string `env:"DB_NAME"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
}
