package sql

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"pkg/errors"
)

type DB struct {
	DB *sqlx.DB
}

func Open(ctx context.Context, driverName string, url string) (*DB, error) {
	db, err := sqlx.Open(driverName, url)
	if err != nil {
		return nil, wrapSQLError(ctx, err)
	}
	return &DB{db}, nil
}

func (s *DB) Close(ctx context.Context) error {
	if err := s.DB.Close(); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *DB) Begin(ctx context.Context) (*Tx, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, wrapSQLError(ctx, err)
	}
	return &Tx{tx}, nil
}

func (s *DB) Ping(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *DB) Unsafe() *DB {
	return &DB{s.DB.Unsafe()}
}

func (s *DB) Select(ctx context.Context, dest any, q sq.Sqlizer) (err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)

	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Выполняем запрос в рамках транзакции
		err = tx.Tx.SelectContext(ctx, dest, query, args...)
	} else {

		// Выполняем запрос
		err = s.DB.SelectContext(ctx, dest, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return wrapSQLError(ctx, err)
	}

	return nil
}

func (s *DB) Get(ctx context.Context, dest any, q sq.Sqlizer) (err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)

	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Выполняем запрос в рамках транзакции
		err = tx.Tx.GetContext(ctx, dest, query, args...)
	} else {

		// Выполняем запрос
		err = s.DB.GetContext(ctx, dest, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return wrapSQLError(ctx, err)
	}

	return nil
}

func (s *DB) Query(ctx context.Context, q sq.Sqlizer) (_ *Rows, err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	rows := &Rows{Rows: nil}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Выполняем запрос в рамках транзакции
		rows.Rows, err = tx.Tx.QueryxContext(ctx, query, args...)
	} else {

		// Выполняем запрос
		rows.Rows, err = s.DB.QueryxContext(ctx, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return nil, wrapSQLError(ctx, err)
	}

	return rows, nil
}

func (s *DB) QueryRow(ctx context.Context, q sq.Sqlizer) (*Row, error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	row := &Row{Row: nil}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Выполняем запрос в рамках транзакции
		row.Row = tx.Tx.QueryRowxContext(ctx, query, args...)
	} else {

		// Выполняем запрос
		row.Row = s.DB.QueryRowxContext(ctx, query, args...)
	}

	return row, nil
}

func (s *DB) Prepare(ctx context.Context, q sq.Sqlizer) (_ *Stmt, err error) {

	// Формируем запрос из билдера
	query, _, err := q.ToSql()
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	var stmt = &Stmt{Stmt: nil}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Подготавливаем запрос в рамках транзакции
		stmt.Stmt, err = tx.Tx.PreparexContext(ctx, query)
	} else {

		// Подготавливаем запрос
		stmt.Stmt, err = s.DB.PreparexContext(ctx, query)
	}

	// Обрабатываем ошибки
	if err != nil {
		return nil, wrapSQLError(ctx, err)
	}

	return stmt, nil
}

func (s *DB) Exec(ctx context.Context, q sq.Sqlizer) (err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)

	if err != nil {
		return errors.InternalServer.Wrap(ctx, err)
	}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Исполняем запрос в рамках транзакции
		_, err = tx.Tx.ExecContext(ctx, query, args...)
	} else {

		// Исполняем запрос
		_, err = s.DB.ExecContext(ctx, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return wrapSQLError(ctx, err)
	}

	return nil
}

func (s *DB) ExecWithLastInsertID(ctx context.Context, q sq.Sqlizer) (id uint32, err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return 0, errors.InternalServer.Wrap(ctx, err)
	}

	query += " RETURNING id"

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)

	if err != nil {
		return 0, errors.InternalServer.Wrap(ctx, err)
	}

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Исполняем запрос в рамках транзакции
		err = tx.Tx.GetContext(ctx, &id, query, args...)
	} else {

		// Исполняем запрос
		err = s.DB.GetContext(ctx, &id, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}

	return id, nil
}

func (s *DB) ExecWithRowsAffected(ctx context.Context, q sq.Sqlizer) (_ uint32, err error) {

	// Формируем запрос из билдера
	query, args, err := q.ToSql()
	if err != nil {
		return 0, errors.InternalServer.Wrap(ctx, err)
	}

	// Заменяем все ? на $1, $2 и т.д.
	query, err = sq.Dollar.ReplacePlaceholders(query)

	if err != nil {
		return 0, errors.InternalServer.Wrap(ctx, err)
	}

	var result sql.Result

	// Извлекаем транзакцию из контекста
	if tx := extractTx(ctx); tx != nil {

		// Исполняем запрос в рамках транзакции
		result, err = tx.Tx.ExecContext(ctx, query, args...)
	} else {

		// Исполняем запрос
		result, err = s.DB.ExecContext(ctx, query, args...)
	}

	// Обрабатываем ошибки
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}

	// Получаем количество затронутых строк
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}

	return uint32(affected), nil
}

func wrapSQLError(ctx context.Context, err error) error {

	thirdPathDepthOption := errors.SkipPreviousCallerOption()

	switch {
	case errors.Is(err, context.Canceled):
		return errors.Timeout.Wrap(ctx, err, thirdPathDepthOption)
	case errors.Is(err, sql.ErrNoRows):
		return errors.NotFound.Wrap(ctx, err, thirdPathDepthOption)
	default:
		return errors.InternalServer.Wrap(ctx, err, thirdPathDepthOption)
	}
}
