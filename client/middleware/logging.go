package middleware

import (
	"context"
	"time"

	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
	tleclientsvc "github.com/thelotter-enterprise/usergo/client/service"
)

// NewLoggingMiddleware ...
func NewLoggingMiddleware(logger tlelogger.Logger) tleclientsvc.ServiceMiddleware {
	return func(next tleclientsvc.Service) tleclientsvc.Service {
		return loggingmw{logger, next}
	}
}

type loggingmw struct {
	logger tlelogger.Logger
	next   tleclientsvc.Service
}

func (mw loggingmw) GetUserByID(ctx context.Context, id int) (response tlehttp.Response) {
	defer func(begin time.Time) {
		logger := mw.logger
		_ = tlelogger.InfoWithContext(
			context.Background(),
			logger,
			"GetUseByID middleware",
			"method", "GetUserByID",
			"input", id,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByID(ctx, id)
}

func (mw loggingmw) GetUserByEmail(ctx context.Context, email string) (response tlehttp.Response) {
	defer func(begin time.Time) {
		logger := mw.logger
		_ = tlelogger.InfoWithContext(
			context.Background(),
			logger,
			"GetUseByEmail middleware",
			"method", "GetUserByEmail",
			"input", email,
			"output", response,
			"err", response.Error,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUserByEmail(ctx, email)
}
