package models

// Config  start parameters for lunch the service
type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabasePG    string `env:"DATABASE_PG_URI"`
	DatabaseClick string `env:"DATABASE_Click_URI"`
}
