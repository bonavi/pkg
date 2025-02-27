package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Stmt struct {
	Stmt *sqlx.Stmt
}

func (s *Stmt) Select(ctx context.Context, dest any, args ...any) error {
	if err := s.Stmt.SelectContext(ctx, dest, args...); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Stmt) Get(ctx context.Context, dest any, args ...any) error {
	if err := s.Stmt.GetContext(ctx, dest, args...); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Stmt) Exec(ctx context.Context, args ...any) error {
	_, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Stmt) ExecWithLastInsertID(ctx context.Context, args ...any) (uint32, error) {
	res, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}
	return uint32(id), nil
}

func (s *Stmt) ExecWithAffectedRows(ctx context.Context, args ...any) (uint32, error) {
	res, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, wrapSQLError(ctx, err)
	}
	return uint32(rows), nil
}

func (s *Stmt) QueryRow(ctx context.Context, args ...any) *Row {
	return &Row{s.Stmt.QueryRowxContext(ctx, args...)}
}

func (s *Stmt) Query(ctx context.Context, args ...any) (*Rows, error) {
	rows, err := s.Stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, wrapSQLError(ctx, err)
	}
	return &Rows{rows}, nil
}

func (s *Stmt) Close(ctx context.Context) error {
	if err := s.Stmt.Close(); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}
