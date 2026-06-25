package telemetry

import (
	"context"
	"fmt"

	"jarvis/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitTracerProvider initializes an OpenTelemetry TracerProvider.
func InitTracerProvider(cfg config.TelemetryConfig) (func(context.Context) error, error) {
	// Minimal tracer provider to avoid requiring OTLP gRPC API during compilation.
	tracerProvider := trace.NewTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

// InitMeterProvider initializes an OpenTelemetry MeterProvider.
func InitMeterProvider(cfg config.TelemetryConfig) (func(context.Context) error, error) {
	// Prometheus exporter for metrics
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", "development"), // TODO: Get from config
		)),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
