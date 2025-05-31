package sqlite

import (
	"context"
	"pkg/sql"
)

type SQLiteConfigEnv struct {
	Path string `env:"SQLITE_PATH"`
}

func NewClientSQLite(ctx context.Context, config SQLiteConfigEnv) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db.Unsafe(), nil
}
