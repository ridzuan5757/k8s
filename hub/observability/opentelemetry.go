package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func initConn() (*grpc.ClientConn, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	config := &tls.Config{
		RootCAs: rootCAs,
	}

	credentials.NewTLS(config)

	conn, err := grpc.NewClient(
		"localhost:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection to collector: %v", err)
	}

	return conn, err
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("hub-monitoring"),
			semconv.ServiceVersion("0.0.1"),
		),
	)
}

func newMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {

	otlpMetricsExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(otlpMetricsExporter,
			metric.WithInterval(10*time.Second))),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	conn, err := initConn()
	if err != nil {
		log.Fatal(err)
	}

	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handlerErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	resource, err := newResource()
	if err != nil {
		handlerErr(err)
		return
	}

	meterProviderShutdown, err := newMeterProvider(ctx, resource, conn)
	if err != nil {
		handlerErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProviderShutdown)

	return
}
