package middleware

import (
	"context"

	metrics "github.com/go-kit/kit/metrics"
	tlelogger "github.com/kolbis/corego/logger"
	tlemetrics "github.com/kolbis/corego/metrics"
	"github.com/kolbis/go-kit-user-ms-example/shared"
	"github.com/kolbis/go-kit-user-ms-example/svc"
)

// NewInstrumentingMiddleware ..
func NewInstrumentingMiddleware(logger tlelogger.Logger, inst tlemetrics.PrometheusInstrumentor) ServiceMiddleware {
	counter := inst.AddPromCounter("user", "getuserbyid", tlemetrics.RequestCount, []string{"method", "error"})
	requestLatency := inst.AddPromSummary("user", "getuserbyid", tlemetrics.LatencyInMili, []string{"method", "error"})

	return func(next svc.Service) svc.Service {
		mw := instrumentingMiddleware{
			next:           next,
			requestCount:   counter,
			requestLatency: requestLatency,
			logger:         logger,
		}
		return mw
	}
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           svc.Service
	logger         tlelogger.Logger
}

func (mw instrumentingMiddleware) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	return mw.next.GetUserByID(ctx, userID)
}

func (mw instrumentingMiddleware) ConsumeLoginCommand(ctx context.Context, userID int) error {
	return mw.next.ConsumeLoginCommand(ctx, userID)
}
