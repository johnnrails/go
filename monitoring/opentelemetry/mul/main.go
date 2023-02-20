package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)


func NewExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func NewResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("mul"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

const name = "mult2"

type App struct {
	r io.Reader
	l *log.Logger
}


func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

func (a *App) Run(ctx context.Context) error {
	for {
		newCtx, span := otel.Tracer(name).Start(ctx, "run")
		n, err := a.Poll(ctx)
		if err != nil {
			span.End()
			return err
		}
		a.Write(newCtx, n)
		span.End()
	}
}

func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("number: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, err
}

func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	a.l.Printf("%d * 2 = %d\n", n, n * 2)

	f := func(ctx context.Context) (uint) {
		_, span := otel.Tracer(name).Start(ctx, "Multiply")
		defer span.End()
		
		return n * 2
	}(ctx)

	a.l.Printf("Multiply(%d) = %d\n", n, f)
}

func main() {
	l := log.New(os.Stdout, "", 0)
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(f)
	}
	defer f.Close()

	// Configure exporter
	exp, err := NewExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	// Configure trace
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(NewResource()),
	)
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()

	// Set the trace to opentelemetry
	otel.SetTracerProvider(traceProvider)

	// Notify the signal channel when has Interrupt (CTRL+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Error channel to see when has errors
	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	
	go func ()  {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <- sigCh:
		l.Print("\ngoodbye")
	case err := <- errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}