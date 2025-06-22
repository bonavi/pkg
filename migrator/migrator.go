package migrator

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"pkg/errors"

	"pkg/log"

	"github.com/pressly/goose/v3"
)

type MigratorConfig struct {
	Conn            *sql.DB            // Подключение к базе данных
	EmbedMigrations embed.FS           // Встроенные файлы миграций
	Dialect         goose.Dialect      // Драйвер
	Dir             string             // Путь к миграциям, так как embedding сохраняет структуру директорий
	Migrations      []*goose.Migration // Миграции
}

type Migrator struct {
	provider *goose.Provider
}

func NewMigrator(config MigratorConfig) (res Migrator, err error) {
	provider, err := goose.NewProvider(config.Dialect, config.Conn, config.EmbedMigrations, goose.WithGoMigrations(config.Migrations...), goose.WithAllowOutofOrder(true))
	if err != nil {
		return res, errors.Default.Wrap(err)
	}

	return Migrator{
		provider: provider,
	}, nil
}

func (m Migrator) Up(ctx context.Context) error {

	result, err := m.provider.Up(ctx)
	if err != nil {
		return errors.Default.Wrap(err)
	}

	for _, r := range result {
		if r.Error != nil {
			return r.Error
		}

		log.Info(fmt.Sprintf("migration %d applied", r.Source.Version))
	}

	return nil
}
