package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound - запись не найдена
	ErrNotFound = errors.New("record not found")
)

func dbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return errors.WithStack(err)
}
