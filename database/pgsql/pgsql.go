package pgsql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:golint

	"pkg/sql"
)

func NewClientPgsql(ctx context.Context, connectionURI string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionURI)
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

func GetConnectionURI(host, user, password, database string) string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", user, password, host, database)
}
