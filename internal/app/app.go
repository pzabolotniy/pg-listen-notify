package app

import (
	"context"
	"net/http"

	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/conf"
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
