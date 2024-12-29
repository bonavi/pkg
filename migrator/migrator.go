package migrator

import (
	"context"
	"embed"

	"pkg/errors"
	"pkg/sql"

	"github.com/pressly/goose/v3"
)

type Migrator interface {
	Up(context.Context) error
	Down(context.Context) error
}

// Dialect is the type of database dialect.
type Dialect string

const (
	DialectClickHouse Dialect = "clickhouse"
	DialectMySQL      Dialect = "mysql"
	DialectPostgres   Dialect = "postgres"
	DialectSQLite3    Dialect = "sqlite3"
	DialectYdB        Dialect = "ydb"
)

type MigratorConfig struct {
	EmbedMigrations embed.FS // Встроенные файлы миграций
	Dialect         Dialect  // Драйвер ex: clickhouse
	Dir             string   // Путь к миграциям, так как embedding сохраняет структуру директорий
}

type migrator struct {
	cfg  MigratorConfig
	conn *sql.DB
}

func NewMigrator(conn *sql.DB, config MigratorConfig) (Migrator, error) {

	goose.SetBaseFS(config.EmbedMigrations)

	goose.SetLogger(newMigratorLogger())

	if err := goose.SetDialect(string(config.Dialect)); err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	return migrator{
		conn: conn,
		cfg:  config,
	}, nil
}

func (mg migrator) Up(ctx context.Context) error {
	if err := goose.UpContext(ctx, mg.conn.DB.DB, mg.cfg.Dir); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	return nil
}

func (mg migrator) Down(ctx context.Context) error {
	if err := goose.DownContext(ctx, mg.conn.DB.DB, mg.cfg.Dir); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	return nil
}
