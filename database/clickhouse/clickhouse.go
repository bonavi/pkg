package clickhouse

import (
	"context"
	"fmt"

	"pkg/database"
	"pkg/sql"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

const (
	maxOpenConns = 10
	maxIdleConns = 5
	maxBlockSize = 10
)

type ClickhouseConfig struct {
	Host         string `env:"CLICKHOUSE_HOST"`
	User         string `env:"CLICKHOUSE_USER"`
	Password     string `env:"CLICKHOUSE_PASSWORD"`
	DatabaseName string `env:"CLICKHOUSE_DATABASE"`
}

func NewClientClickhouse(config ClickhouseConfig) (*sql.DB, error) {

	db, err := sql.Open("clickhouse", config.getConnectionURI())

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.ConnectionTimeout)
	defer cancel()

	if err = db.Ping(ctx); err != nil {

		return nil, err
	}

	return db, err
}

func (c *ClickhouseConfig) getConnectionURI() string {
	return fmt.Sprintf("clickhouse://%s:%s@%s/%s?dial_timeout=%v&max_open_conns=%v&max_idle_conns=%v&conn_max_lifetime=%v&max_block_size=%v",
		c.User,
		c.Password,
		c.Host,
		c.DatabaseName,
		"1s",         // dial_timeout
		maxOpenConns, // max_open_conns
		maxIdleConns, // max_idle_conns
		"1h",         // conn_max_lifetime
		maxBlockSize, // max_block_size
	)

}
