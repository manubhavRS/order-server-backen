package models

type DatabaseConfig struct {
	Host         string `env:"host"`
	Port         string `env:"port"`
	DatabaseName string `env:"database"`
	User         string `env:"user"`
	Password     string `env:"password"`
	ServerPort   string `env:"serverPort"`
}
