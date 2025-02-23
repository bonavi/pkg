package pgsql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:golint

	"pkg/sql"
)

type PostgreSQLConfig struct {
	Host     string `env:"PGSQL_HOST,required"`
	User     string `env:"PGSQL_USER,required"`
	Password string `env:"PGSQL_PASSWORD,required"`
	Database string `env:"PGSQL_DATABASE,required"`
}

func (c *PostgreSQLConfig) getURL() string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", c.User, c.Password, c.Host, c.Database)
}

func NewClientPgsql(ctx context.Context, conf PostgreSQLConfig) (*sql.DB, error) {
	db, err := sql.Open(ctx, "pgx", conf.getURL())
	if err != nil {
		return nil, err
	}

	// 	TODO вынести в конфиги
	const (
		MaxIdleConns = 50
		MaxOpenConns = 500
	)

	db.DB.SetMaxOpenConns(MaxOpenConns)
	db.DB.SetMaxIdleConns(MaxIdleConns)

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db.Unsafe(), nil
}
