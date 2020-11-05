package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tlesd "github.com/thelotter-enterprise/corego/servicediscovery"
	"github.com/thelotter-enterprise/usergo/client"
)

func TestClientIntegration(t *testing.T) {
	serviceName := "test"
	logger := tlelogger.NewNopLogger()
	ctx := context.Background()
	id := 2

	consulServiceDiscoverator := makeConsulServiceDiscovery(logger)
	dnsServiceDiscoverator := makeDNSServiceDiscovery(logger)
	c := client.NewServiceClientWithDefaults(logger, consulServiceDiscoverator, dnsServiceDiscoverator, serviceName)

	response := c.GetUserByID(ctx, id)

	if response.Data == nil {
		t.Fail()
	}
}

func makeLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", ":8080", "caller", log.DefaultCaller)

	return logger
}

func makeConsulServiceDiscovery(logger tlelogger.Logger) *tlesd.ConsulServiceDiscovery {
	consulAddress := "localhost:8500"
	sd := tlesd.NewConsulServiceDiscovery(logger, consulAddress)
	return &sd
}

func makeDNSServiceDiscovery(logger tlelogger.Logger) *tlesd.DNSServiceDiscovery {
	sd := tlesd.NewDNSServiceDiscovery(logger)
	return &sd
}
