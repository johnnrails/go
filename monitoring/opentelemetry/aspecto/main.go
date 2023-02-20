package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func startWebServer() {
  r := gin.Default()
  //gin OpenTelemetry instrumentation
  r.Use(otelgin.Middleware("aspecto-service"))
  r.GET("/todo", func(c *gin.Context) {
      results := map[string]interface{}{
        "hello": "hello",
      }
      c.JSON(http.StatusOK, results)
  })
  _ = r.Run(":8080")
}

func JaegerTracerProvider() (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
  if err != nil {
    return nil, err
  }
  tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exp),
    sdktrace.WithResource(resource.NewWithAttributes(
      semconv.SchemaURL,
      semconv.ServiceNameKey.String("aspecto-service"),
      semconv.DeploymentEnvironmentKey.String("development"),
    )),
  )
  return tp, nil
}

func main() {
	ctx := context.Background()

	traceExporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("collector.aspecto.io"),
		otlptracehttp.WithHeaders(map[string]string{"Authorization": "23b211a6-01fb-4dad-b8e8-c8ca40177bcf"}),
	)

	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("aspecto-service"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(tp)

	// jTp, _ := JaegerTracerProvider()
	// otel.SetTracerProvider(jTp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	startWebServer()
}
