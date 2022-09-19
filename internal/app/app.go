package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/listener"
)

func StartWebAPI(ctx context.Context, router http.Handler, webAPI *conf.WebAPI) error {
	logger := logging.FromContext(ctx)
	logger.WithField("listen", webAPI.Listen).Trace("listen addr")
	if err := http.ListenAndServe(webAPI.Listen, router); err != nil {
		logger.WithError(err).WithField("listen", webAPI.Listen).Error("listen failed")

		return err
	}

	return nil
}

func StartListener(ctx context.Context, dbConn *pgxpool.Pool, config *conf.Events) error {
	logger := logging.FromContext(ctx)
	logger.Trace("starting pg events listener")
	if err := listener.Serve(ctx, dbConn, config); err != nil {
		logger.WithError(err).Error("listener serve failed")

		return fmt.Errorf("listener serve failed: %w", err)
	}

	return nil
}
