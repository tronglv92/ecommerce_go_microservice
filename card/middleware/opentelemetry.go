package middleware

import (
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type config struct {
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

const (
	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// SemVersion is the semantic version to be supplied to tracer/meter creation.
func SemVersion() string {
	return "semver:" + Version()
}
func Version() string {
	return "0.32.0"
	// This string is updated by the pre_release.sh script during release
}
