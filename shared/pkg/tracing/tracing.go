package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

type Config struct {
	ServiceName  string
	Environment  string
	JaegerAgent  string
	CollectorURL string
	Enabled      bool
}

func InitTracer(cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		// Return a no-op tracer provider
		tp := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		Tracer = tp.Tracer(cfg.ServiceName)
		return tp.Shutdown, nil
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create Jaeger exporter with correct API
	var exp sdktrace.SpanExporter
	if cfg.CollectorURL != "" {
		return nil, fmt.Errorf("OTLP exporter not implemented in this version")
	} else {
		// Use Jaeger exporter (for development)
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerAgent)))
		if err != nil {
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	Tracer = tp.Tracer(cfg.ServiceName)

	return tp.Shutdown, nil
}

// ... rest of the file remains the same
