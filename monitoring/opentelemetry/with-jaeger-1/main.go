package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service = "go-opentelemetry-with-jeager"
	environment = "development"
	id = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// create the jeager exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	// Telemetry data can be crucial to solving issues with a service. The catch is, you need a way to identify what service, 
	// or even what service instance, that data is coming from. OpenTelemetry uses a Resource to represent the entity producing telemetry.
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		// Record information about this app in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)

	return tracerProvider, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

  // Register our TracerProvider as the global so any imported
  // instrumentation in the future will default to using it.
  otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	 // Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
				log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("main")
	
	ctx, span := tr.Start(ctx, "hello")
	defer span.End()

	helloHandler := func (w http.ResponseWriter, r *http.Request)  {
		tr := otel.Tracer("hello-handler")
		_, span := tr.Start(ctx, "hello")
		span.SetAttributes(attribute.Key("key").String("value"))
		defer span.End()

		name := os.Getenv("MY_NAME")
		fmt.Fprintf(w, "hello %q!", name)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	http.Handle("/", otelHandler)

	log.Println("Listening on localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}