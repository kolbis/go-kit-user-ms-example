package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kithttp "github.com/go-kit/kit/transport/http"
	tlectx "github.com/kolbis/corego/context"
	tlelogger "github.com/kolbis/corego/logger"
	tlemetrics "github.com/kolbis/corego/metrics"
	tletracer "github.com/kolbis/corego/tracer"
	tlerabbitmq "github.com/kolbis/corego/transport/rabbitmq"
	tleimpl "github.com/kolbis/go-kit-user-ms-example/svc/implementation"
	svcmw "github.com/kolbis/go-kit-user-ms-example/svc/middleware"
	svctrans "github.com/kolbis/go-kit-user-ms-example/svc/transport"
	svcamqp "github.com/kolbis/go-kit-user-ms-example/svc/transport/amqp"
	svchttp "github.com/kolbis/go-kit-user-ms-example/svc/transport/http"
)

func main() {
	var (
		serviceName      string                    = "user"
		hostAddress      string                    = "localhost:8080"
		zipkinURL        string                    = "http://localhost:9411/api/v2/spans"
		rabbitMQUsername string                    = "thelotter"
		rabbitMQPwd      string                    = "Dhvbuo1"
		rabbitMQHost     string                    = "int-k8s1"
		rabbitMQVhost    string                    = "thelotter"
		rabbitMQPort     int                       = 32672
		env              string                    = "dev"
		logLevel         tlelogger.AtomicLevelName = tlelogger.DebugLogLevel
		ctx              context.Context           = tlectx.Root()
	)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// infra services
	logger, _ := tlelogger.NewLogger(env, "Fiel", logLevel)
	tracer := tletracer.NewTracer(serviceName, hostAddress, zipkinURL)
	promMetricsInst := tlemetrics.NewPrometheusInstrumentor(serviceName)

	// implementation
	repo := tleimpl.NewRepository()
	service := tleimpl.NewService(logger, tracer, repo)

	// middlewares
	service = svcmw.NewLoggingMiddleware(logger)(service)                        // Hook up the logging middleware
	service = svcmw.NewInstrumentingMiddleware(logger, promMetricsInst)(service) // Hook up the inst middleware

	// making all types of endpoints
	endpoints := svctrans.MakeEndpoints(service)

	// http
	handler := svchttp.NewTransport(ctx, endpoints, make([]kithttp.ServerOption, 0), logger)
	go func() {
		server := &http.Server{
			Addr:    hostAddress,
			Handler: handler,
		}
		tlelogger.InfoWithContext(ctx, logger, fmt.Sprintf("listening for http calls on %s", hostAddress))
		errs <- server.ListenAndServe()
		done <- true
	}()

	// setting up RabbitMQ server
	connInfo := tlerabbitmq.NewConnectionInfo(rabbitMQHost, rabbitMQPort, rabbitMQUsername, rabbitMQPwd, rabbitMQVhost)
	conn := tlerabbitmq.NewConnectionManager(connInfo)
	subscribers := svcamqp.NewTransport(endpoints, logger, &conn)
	publisher := tlerabbitmq.NewPublisher(&conn)
	client := tlerabbitmq.NewClient(&conn, logger, &publisher, subscribers)
	amqpServer := tlerabbitmq.NewServer(logger, tracer, &client, &conn)

	go func() {
		tlelogger.InfoWithContext(ctx, logger, fmt.Sprintf("listening for amqp messages"))
		err := amqpServer.Run(ctx)
		if err != nil {
			errs <- err
			fmt.Println(err)
			done <- true
		}
	}()

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
