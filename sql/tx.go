package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Tx struct {
	Tx *sqlx.Tx
}

func (s *Tx) Commit(ctx context.Context) error {
	if err := s.Tx.Commit(); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}

func (s *Tx) Rollback(ctx context.Context) error {
	if err := s.Tx.Rollback(); err != nil {
		return wrapSQLError(ctx, err)
	}
	return nil
}
