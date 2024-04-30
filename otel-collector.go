package main

import (
	"context"
	"fmt"
	otelmetric "go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Initializes an OTLP exporter, and configures the corresponding metrics provider.
func main() {

	ctx := context.Background()
	res, err := newResource()
	if err != nil {
		panic(err)
	}

	mp, err := newMeterProvider(res)
	if err != nil {
		panic(err)
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	meter := mp.Meter("test")

	counter, err := meter.Int64Counter("test.counter",
		otelmetric.WithUnit("1"),
		otelmetric.WithDescription("test counter"))
	if err != nil {
		panic(err)
	}

	otel.SetMeterProvider(mp)

	for i := 0; i < 10; i++ {
		counter.Add(ctx, 1)
	}

}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("Test"),
			semconv.ServiceVersion("0.0.1"),
		))
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// using NodePort service to connect to the otel collector running in k8s
	conn, err := grpc.NewClient("172.18.0.2:30080",
		// TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a metrics exporter
	metricExporter, err := otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(1*time.Second))),
	)
	return meterProvider, nil
}
