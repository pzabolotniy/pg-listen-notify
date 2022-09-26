package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/listener"
)

func StartWebAPI(ctx context.Context, router http.Handler, webAPI *conf.WebAPI) error {
	logger := logging.FromContext(ctx)

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

func StartListener(ctx context.Context, dbConn *pgxpool.Pool, config *conf.Events) error {
	logger := logging.FromContext(ctx)

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

	serveErrCh := make(chan error, 1)
	logger.Trace("starting pg events listener")
	server := listener.NewServer()
	go func() {
		if err = server.Serve(ctx, dbConn, config); err != nil && errors.Is(err, listener.ErrGracefulShutdown) {
			logger.WithError(err).Error("listener serve failed")

			serveErrCh <- fmt.Errorf("listener serve failed: %w", err)
		}
	}()

	for {
		select {
		case s := <-gracefulShutdown:
			logger.WithField("signal", s.String()).Trace("signal caught. terminating listener")
			server.Shutdown()

			return nil
		case otherErr := <-serveErrCh:
			logger.WithError(otherErr).Trace("terminating listener")

			return nil
		}
	}
}
