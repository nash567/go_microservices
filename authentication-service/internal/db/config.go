package db

import "fmt"

type Config struct {
	Name      string
	Host      string
	Port      uint16
	User      string
	Password  string
	MaxConns  uint
	IdleConns uint
	Driver    string
}

// "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

func (cfg *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
