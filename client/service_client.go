package client

import (
	"context"

	"github.com/gorilla/mux"
	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tlemetrics "github.com/thelotter-enterprise/corego/metrics"
	tleratelimit "github.com/thelotter-enterprise/corego/ratelimit"
	tlesd "github.com/thelotter-enterprise/corego/servicediscovery"
	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
	tleclientmw "github.com/thelotter-enterprise/usergo/client/middleware"
	tleclientsvc "github.com/thelotter-enterprise/usergo/client/service"
)

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      tlelogger.Logger
	Consul      *tlesd.ConsulServiceDiscovery
	ServiceName string
	DNS         *tlesd.DNSServiceDiscovery
	Limiter     tleratelimit.RateLimiterConfig
	Inst        tlemetrics.PrometheusInstrumentor
	Router      *mux.Router
}

// NewServiceClientWithDefaults with defaults
func NewServiceClientWithDefaults(logger tlelogger.Logger, consul *tlesd.ConsulServiceDiscovery, dns *tlesd.DNSServiceDiscovery, serviceName string) ServiceClient {
	return NewServiceClient(
		logger,
		consul,
		dns,
		tleratelimit.NewDefaultRateLimiterConfig(),
		tlemetrics.NewPrometheusInstrumentor(serviceName),
		mux.NewRouter(),
		serviceName,
	)
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient(logger tlelogger.Logger, consul *tlesd.ConsulServiceDiscovery, dns *tlesd.DNSServiceDiscovery, limiter tleratelimit.RateLimiterConfig, inst tlemetrics.PrometheusInstrumentor, router *mux.Router, serviceName string) ServiceClient {
	client := ServiceClient{
		Logger:      logger,
		Consul:      consul,
		DNS:         dns,
		ServiceName: serviceName,
		Limiter:     limiter,
		Inst:        inst,
		Router:      router,
	}
	return client
}

// GetUserByID , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client ServiceClient) GetUserByID(ctx context.Context, id int) tlehttp.Response {
	var service tleclientsvc.Service

	instMiddleware := tleclientmw.NewInstrumentingMiddleware(client.Inst)
	logMiddleware := tleclientmw.NewLoggingMiddleware(client.Logger)
	serviceMiddleware := tleclientmw.NewTransportMiddleware(ctx, tleclientmw.GetUserByIDFactory(ctx, id), client.Consul, client.Logger)

	service = serviceMiddleware(service)
	service = logMiddleware(service)
	service = instMiddleware(service)

	res := service.GetUserByID(ctx, id)
	return res
}

// GetUserByEmail , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client ServiceClient) GetUserByEmail(ctx context.Context, email string) tlehttp.Response {
	return tlehttp.Response{}
}
