package config

type Config struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"9000"`

	DB string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/schoolmat?sslmode=disable"`
}
