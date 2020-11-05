package middleware

import (
	"context"
	"fmt"
	"time"

	metrics "github.com/go-kit/kit/metrics"
	tleinst "github.com/kolbis/corego/metrics"
	tlehttp "github.com/kolbis/corego/transport/http"
	tleclientsvc "github.com/kolbis/go-kit-user-ms-example/client/service"
)

// NewInstrumentingMiddleware ...
func NewInstrumentingMiddleware(inst tleinst.PrometheusInstrumentor) tleclientsvc.ServiceMiddleware {
	counter := inst.AddPromCounter("user", "getuserbyid", tleinst.RequestCount, []string{"method", "error"})
	requestLatency := inst.AddPromSummary("user", "getuserbyid", tleinst.LatencyInMili, []string{"method", "error"})

	return func(next tleclientsvc.Service) tleclientsvc.Service {
		mw := instmw{
			next:           next,
			requestCount:   counter,
			requestLatency: requestLatency,
		}
		return mw
	}
}

type instmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           tleclientsvc.Service
}

func (mw instmw) GetUserByID(ctx context.Context, id int) (response tlehttp.Response) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByID", "error", fmt.Sprint(response.Error != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response = mw.next.GetUserByID(ctx, id)
	return response
}

func (mw instmw) GetUserByEmail(ctx context.Context, email string) (response tlehttp.Response) {
	return mw.next.GetUserByEmail(ctx, email)
}
