package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/db"
	"github.com/pzabolotniy/listen-notify/internal/listener"
)

func StartWebAPI(ctx context.Context, logger logging.Logger, router http.Handler, webAPI *conf.WebAPI) error {
	logger = logging.FromContext(ctx, logger)

	tracingProvider, err := initJaegerTracing(logger)
	if err != nil {
		logger.WithError(err).Error("init jaeger tracing failed")

		return err
	}
	defer func() {
		if stopErr := tracingProvider.Shutdown(ctx); stopErr != nil {
			logger.WithError(stopErr).Error("shutting down tracer provider failed")
		}
	}()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM, syscall.SIGINT)

	httpServer := &http.Server{
		Addr:    webAPI.Listen,
		Handler: router,
	}

	serveErrCh := make(chan error, 1)
	logger.WithField("listen", webAPI.Listen).Trace("listen addr")
	go func() {
		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).WithField("listen", webAPI.Listen).Error("listen failed")
			serveErrCh <- err
		}
	}()

	for {
		select {
		case s := <-gracefulShutdown:
			logger.WithField("signal", s.String()).Trace("signal caught. terminating webapi")
			if shutdownErr := httpServer.Shutdown(ctx); shutdownErr != nil {
				logger.WithError(shutdownErr).Error("terminate webapi failed")
			}

			return nil
		case otherErr := <-serveErrCh:
			logger.WithError(otherErr).Trace("terminating webapi")
		}
	}
}

func StartListener(ctx context.Context, logger logging.Logger, dbService *db.DBService, config *conf.Events) error {
	logger = logging.FromContext(ctx, logger)

	tracingProvider, err := initJaegerTracing(logger)
	if err != nil {
		logger.WithError(err).Error("init jaeger tracing failed")

		return err
	}
	defer func() {
		if stopErr := tracingProvider.Shutdown(ctx); stopErr != nil {
			logger.WithError(stopErr).Error("shutting down tracer provider failed")
		}
	}()

	logger.Trace("starting pg events listener")
	if err = listener.Serve(ctx, logger, dbService, config); err != nil {
		logger.WithError(err).Error("listener serve failed")

		return fmt.Errorf("listener serve failed: %w", err)
	}

	return nil
}
