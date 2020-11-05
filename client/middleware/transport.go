package middleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	tleloadbalancer "github.com/thelotter-enterprise/corego/loadbalancer"
	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tlesd "github.com/thelotter-enterprise/corego/servicediscovery"
	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
	tleclientsvc "github.com/thelotter-enterprise/usergo/client/service"
)

// transportmw is used to make the actual call on the endpoint
// Either by calling an http endpoint or rabbitmq
type transportmw struct {
	// Next is a the service instance
	// We need to use Next, since it is used to satisfy the middleware pattern
	// Each middleware is responbsible for a single API, yet, due to the service interface,
	// it need to implement all the service interface APIs. To support it, we use Next to obstract the implementation
	Next interface{}

	// This is the current API which we plan to support in the service interface contract
	This endpoint.Endpoint
}

// NewTransportMiddleware will create a new instance of TransportMiddleware
// It will configure consul to look for user service
// (it is possible to use DNS and not consul)
func NewTransportMiddleware(ctx context.Context, fac sd.Factory, consul *tlesd.ConsulServiceDiscovery, logger tlelogger.Logger) tleclientsvc.ServiceMiddleware {
	consulInstancer, _ := consul.ConsulInstance(ctx, "user", []string{}, true)
	//dnsInstancer := proxy.dns.DNSInstance("user")
	endpointer := sd.NewEndpointer(consulInstancer, fac, logger)
	lb := tleloadbalancer.NewDynamicLoadBalancer(endpointer)
	retry := lb.DefaultRoundRobinWithRetryEndpoint(ctx)

	return func(next tleclientsvc.Service) tleclientsvc.Service {
		return transportmw{Next: next, This: retry}
	}
}

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
// We do this to satisfy the service interface
func (proxymw transportmw) GetUserByID(ctx context.Context, id int) tlehttp.Response {
	response := tlehttp.Execute(ctx, id, proxymw.This)
	return response
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (proxymw transportmw) GetUserByEmail(ctx context.Context, email string) tlehttp.Response {
	svc := proxymw.Next.(tleclientsvc.Service)
	return svc.GetUserByEmail(ctx, email)
}
