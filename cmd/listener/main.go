package main

import (
	"context"
	"fmt"

	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/app"
	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/db"
	"github.com/pzabolotniy/listen-notify/internal/migration"
)

func main() {
	appConf, err := conf.GetConfig()
	if err != nil {
		fmt.Println(err)

		return
	}
	ctx := context.Background()
	logger := logging.GetLogger()
	ctx = logging.WithContext(ctx, logger)

	dbConn, err := db.Connect(ctx, appConf.DB)
	if err != nil {
		logger.WithError(err).Error("db connect failed. exiting.")

		return
	}
	defer db.Disconnect(dbConn)

	err = migration.MigrateUp(ctx, dbConn, appConf.DB)
	if err != nil {
		logger.WithError(err).Error("migration failed")

		return
	}

	startErr := app.StartListener(ctx, dbConn, appConf.Events)
	if startErr != nil {
		logger.WithError(startErr).Error("start events listener failed")
	}
}
