package pgsql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"pkg/sql"
)

type PgsqlConfigEnv struct {
	Host     string `env:"PGSQL_HOST"`
	User     string `env:"PGSQL_USER"`
	Password string `env:"PGSQL_PASSWORD"`
	Database string `env:"PGSQL_DATABASE"`
}

func NewClientPgsql(ctx context.Context, conf PgsqlConfigEnv) (*sql.DB, error) {
	db, err := sql.Open("pgx", conf.GetConnectionURI())
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

func (c *PgsqlConfigEnv) GetConnectionURI() string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", c.User, c.Password, c.Host, c.Database)
}
