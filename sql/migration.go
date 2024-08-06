package sql

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func SetMigration(ctx context.Context, driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driver); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return err
	}

	return nil
}
