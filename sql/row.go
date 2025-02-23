package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Row struct {
	Row *sqlx.Row
}

func (s *Row) Scan(ctx context.Context, dest ...any) error {
	if err := s.Row.Scan(dest...); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}
