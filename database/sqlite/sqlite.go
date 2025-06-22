package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"pkg/errors"
	"pkg/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConfigEnv struct {
	Path string `env:"SQLITE_PATH"`
}

func NewClientSQLite(ctx context.Context, config SQLiteConfigEnv) (*sql.DB, error) {

	// Извлекаем директорию из полного пути
	dir := filepath.Dir(config.Path)

	// Проверяем, существует ли директория для файла базы данных
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		// Создаем все промежуточные директории
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, errors.Default.Wrap(err)
		}
	}

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
