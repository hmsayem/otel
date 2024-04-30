package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {
	ctx := context.Background()
	// Create the exporter - let's use a stdout exporter
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())

	if err != nil {
		panic(err)
	}

	// Create the resource to be traced
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Test"),
			semconv.ServiceVersion("v0.0.1"),
		),
	)
	if err != nil {
		panic(err)
	}

	// Configure the meter provider
	meterProvider := api.NewMeterProvider(
		api.WithResource(res),
		api.WithReader(api.NewPeriodicReader(metricExporter,
			api.WithInterval(1*time.Second))),
	)

	defer func() { _ = meterProvider.Shutdown(ctx) }()
	otel.SetMeterProvider(meterProvider)

	meter := otel.Meter("test-meter")
	// Define a counter metric
	counter, err := meter.Int64Counter("test.counter",
		metric.WithDescription("test counter"))
	if err != nil {
		panic(err)
	}

	counter.Add(context.Background(), 1)
	counter.Add(context.Background(), 1)
}
