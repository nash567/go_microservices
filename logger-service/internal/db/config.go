package db

type Config struct {
	Name     string
	Host     string
	Port     uint16
	UserName string
	Password string
}

// mongodb://mongo:27017
// "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
