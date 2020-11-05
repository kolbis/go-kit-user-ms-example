package middleware

import (
	"context"
	"time"

	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tleutils "github.com/thelotter-enterprise/corego/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc"
)

// NewLoggingMiddleware ... ..
func NewLoggingMiddleware(logger tlelogger.Logger) ServiceMiddleware {
	return func(next svc.Service) svc.Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger tlelogger.Logger
	next   svc.Service
}

func (mw loggingMiddleware) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	dt := tleutils.DateTime{}

	defer func(begin time.Time) {
		logger := mw.logger
		tlelogger.InfoWithContext(
			ctx,
			logger,
			"GetUserByID",
			"method", "GetUserByID",
			"took", time.Since(begin),
		)
	}(dt.Now())

	return mw.next.GetUserByID(ctx, userID)
}

func (mw loggingMiddleware) ConsumeLoginCommand(ctx context.Context, userID int) error {
	dt := tleutils.DateTime{}
	defer func(begin time.Time) {
		logger := mw.logger
		_ = tlelogger.InfoWithContext(
			ctx,
			logger,
			"ConsumeLoginCommand",
			"method", "ConsumeLoginCommand",
			"took", time.Since(begin),
		)
	}(dt.Now())

	return mw.next.ConsumeLoginCommand(ctx, userID)
}
