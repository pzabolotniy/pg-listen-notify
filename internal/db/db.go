package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
)

func Connect(ctx context.Context, dbConf *conf.DB) (*pgx.Conn, error) {
	logger := logging.FromContext(ctx)
	connString := dbConf.ConnString
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		logger.WithError(err).Error("connect failed")
		return nil, err
	}
	return conn, nil
}

func Disconnect(ctx context.Context, dbConn *pgx.Conn) error {
	return dbConn.Close(ctx)
}

func NativeDriver(pgConn *pgx.Conn) (*sql.DB, error) {
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
