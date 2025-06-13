package pgsql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:golint

	"pkg/sql"
)

type PgsqlConfigEnv struct {
	Host     string `env:"PGSQL_HOST"`
	Database string `env:"PGSQL_DATABASE"`
	User     string `env:"PGSQL_USER"`
	Password string `env:"PGSQL_PASSWORD"`
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

func (p PgsqlConfigEnv) GetConnectionURI() string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", p.User, p.Password, p.Host, p.Database)
}
