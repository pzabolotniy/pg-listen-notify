package migration

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pzabolotniy/logging/pkg/logging"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/db"
)

func MigrateUp(ctx context.Context, pgConn *pgx.Conn, migrationConf *conf.DB) error {
	logger := logging.FromContext(ctx)
	migrations := &migrate.FileMigrationSource{
		Dir: migrationConf.MigrationDir,
	}

	dbConn, err := db.NativeDriver(pgConn)
	if err != nil {
		logger.WithError(err).Error("get native driver failed")

		return err
	}
	defer func() {
		if closeErr := db.DisconnectNativeDriver(dbConn); closeErr != nil {
			logger.WithError(closeErr).Error("disconnect failed")
		}
	}()

	migrate.SetTable(migrationConf.MigrationTable)
	migrationsApplied, err := migrate.Exec(dbConn, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.WithError(err).Error("migration failed")

		return err
	}
	logger.WithField("migrations_applied", migrationsApplied).Trace("migration succeeded")

	return nil
}
