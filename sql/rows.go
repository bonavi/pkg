package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Rows struct {
	*sqlx.Rows
}

func (s *Rows) MapScan(ctx context.Context, dest map[string]any) error {
	if err := s.Rows.MapScan(dest); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Rows) StructScan(ctx context.Context, dest any) error {
	if err := s.Rows.StructScan(dest); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Rows) Close(ctx context.Context) error {
	if err := s.Rows.Close(); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}
