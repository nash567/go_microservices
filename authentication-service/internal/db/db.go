package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	maxAllowedPacket = 104857600 // 100MB
	txIsolation      = "READ-COMMITTED"
)

func NewConnection(config *Config) (*sql.DB, error) {
	fmt.Println("sql string", config.DSN(), config.Driver)
	conn, err := sql.Open(config.Driver, config.DSN())
	if err != nil {
		fmt.Println("err.....................", err)
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}
	// conn.SetMaxOpenConns(int(config.MaxConns))
	// conn.SetMaxIdleConns(int(config.IdleConns))

	return conn, nil
}

// Verify ensures the connection is available for use.
func Verify(conn *sql.DB) error {
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
