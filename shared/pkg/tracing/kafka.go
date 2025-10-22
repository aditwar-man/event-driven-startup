package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func InjectKafkaTrace(ctx context.Context, headers map[string]string) {
	carrier := propagation.MapCarrier(headers)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func ExtractKafkaTrace(ctx context.Context, headers map[string]string) context.Context {
	carrier := propagation.MapCarrier(headers)
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

func StartKafkaConsumerSpan(ctx context.Context, topic string, partition int, offset int64) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("kafka.consume %s", topic)
	ctx, span := Tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.operation", "receive"),
			attribute.Int("messaging.kafka.partition", partition),
			attribute.Int64("messaging.kafka.offset", offset),
		),
	)
	return ctx, span
}

func StartKafkaProducerSpan(ctx context.Context, topic string, key string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("kafka.produce %s", topic)
	ctx, span := Tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.operation", "send"),
			attribute.String("messaging.kafka.message_key", key),
		),
	)
	return ctx, span
}
