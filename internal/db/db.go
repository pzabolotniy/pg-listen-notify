package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
)

func Connect(ctx context.Context, dbConf *conf.DB) (*pgxpool.Pool, error) {
	logger := logging.FromContext(ctx)
	connString := dbConf.ConnString
	parsedConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.
			WithError(err).
			WithField("conn_string", connString).
			Error("parse connection string failed")

		return nil, fmt.Errorf("parse connection string failed: %w", err)
	}
	parsedConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	conn, err := pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		logger.WithError(err).Error("connect failed")

		return nil, fmt.Errorf("connect failed: %w", err)
	}

	conn.Config().MaxConns = dbConf.MaxOpenConns
	conn.Config().MaxConnIdleTime = dbConf.ConnMaxLifetime

	return conn, nil
}

type Closer interface {
	Close()
}

func Disconnect(dbConn Closer) {
	dbConn.Close()
}

type Configurer interface {
	Config() *pgxpool.Config
}

func NativeDriver(pgConn Configurer) (*sql.DB, error) {
	conn, err := sql.Open("pgx", pgConn.Config().ConnString())
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func DisconnectNativeDriver(dbConn *sql.DB) error {
	return dbConn.Close()
}
