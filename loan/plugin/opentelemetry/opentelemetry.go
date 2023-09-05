package opentelemetry

import (
	"context"
	"flag"

	"github.com/tronglv92/ecommerce_go_common/logger"
	"github.com/tronglv92/loans/common"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type opentelemetry struct {
	logger logger.Logger

	processName       string
	tracer            trace.Tracer
	sampleTraceRating float64

	stdTracingEnabled bool
	tp                *sdktrace.TracerProvider
}

func NewJaeger(processName string) *opentelemetry {
	return &opentelemetry{
		processName: processName,
	}
}

func (j *opentelemetry) Name() string {
	return common.PluginOpenTelemetry
}

func (j *opentelemetry) GetPrefix() string {
	return j.Name()
}

func (j *opentelemetry) Get() interface{} {
	return j.tracer
}

func (j *opentelemetry) InitFlags() {
	flag.Float64Var(
		&j.sampleTraceRating,
		"jaeger-trace-sample-rate",
		1.0,
		"sample rating for remote tracing from OpenSensus: 0.0 -> 1.0 (default is 1.0)",
	)

	flag.BoolVar(
		&j.stdTracingEnabled,
		"jaeger-std-enabled",
		false,
		"enable tracing export to std (default is false)",
	)
}

func (j *opentelemetry) Configure() error {
	j.logger = logger.GetCurrent().GetLogger(j.Name())
	// client := otlptracegrpc.NewClient(
	// 	otlptracegrpc.WithInsecure(),
	// )
	// exporter, err := otlptrace.New(context.Background(), client)
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))

	if err != nil {
		j.logger.Errorf("creating OTLP trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(j.sampleTraceRating))),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(j.processName)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	j.tp = tp
	j.tracer = tp.Tracer(j.processName)
	j.logger.Info("Run jeager success")
	return nil
}
func newResource(service string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(service),
		semconv.ServiceVersion("0.0.1"),
	)
}
func (j *opentelemetry) Run() error {
	if err := j.Configure(); err != nil {
		return err
	}
	return nil
}

func (j *opentelemetry) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		err := j.tp.Shutdown(context.Background())
		if err != nil {
			j.logger.Errorf("shuttown tracking provider err: %w", err)
		}
		c <- true
		j.logger.Infoln("Stopped")
	}()
	return c
}

// func (j *jaeger) isEnabled() bool {
// 	return j.agentURI != ""
// }
