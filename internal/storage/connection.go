package storage

import (
	"context"
	"database/sql"
	"schoolmat/internal/storage/db"
)

type Connection struct {
	dsn string
	db  *sql.DB

	Q *db.Queries
}

func NewConnection(ctx context.Context, databaseDSN string) (Connection, error) {
	dbase, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return Connection{}, err
	}

	err = dbase.PingContext(ctx)
	if err != nil {
		return Connection{}, err
	}

	queries := db.New(dbase)

	connection := Connection{
		dsn: databaseDSN,
		db:  dbase,
		Q:   queries,
	}
	return connection, nil
}

type QuerierHandler func(context.Context, *db.Queries) error

func (c *Connection) QueryWithTx(ctx context.Context, txopts sql.TxOptions, f QuerierHandler) error {
	tx, err := c.db.BeginTx(ctx, &txopts)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := c.Q.WithTx(tx)
	err = f(ctx, queries)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
