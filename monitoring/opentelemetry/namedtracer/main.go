package main

import (
	"context"
	"fmt"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	fooKey     = attribute.Key("ex.com/foo")
	barKey     = attribute.Key("ex.com/bar")
	anotherKey = attribute.Key("ex.com/another")
)

func SubOperation(ctx context.Context) error {
	tr := otel.Tracer("opentelemetry/namedtracer/suboperation")
	var span trace.Span
	_, span = tr.Start(ctx, "Sub operation...")
	defer span.End()
	span.SetAttributes(attribute.Key("ex.com/lemons").String("five"))
	span.AddEvent("Sub span event")
	return nil
}

func main() {
	stdr.SetVerbosity(5)

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		err = fmt.Errorf("failed to initialize stdouttrace: %w", err)
		panic(err)
	}
	
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("opentelemetry/namedtracer/main")
	ctx := context.Background()
	defer func(){ _ = tracerProvider.Shutdown(ctx) }()

	m0, _ := baggage.NewMember(string(fooKey), "foo1")
	m1, _ := baggage.NewMember(string(barKey), "bar1")
	b, _ := baggage.New(m0, m1)
	ctx = baggage.ContextWithBaggage(ctx, b)

	var span trace.Span
	ctx, span = tracer.Start(ctx, "operation")
	defer span.End()

	span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
	span.SetAttributes(anotherKey.String("yes"))

	if err := SubOperation(ctx); err != nil {
		panic(err)
	}
}