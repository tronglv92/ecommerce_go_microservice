package middleware

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tronglv92/cards/common"
	"github.com/tronglv92/ecommerce_go_common/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// func Recover(ac appctx.AppContext) gin.HandlerFunc {
func Recover() gin.HandlerFunc {
	cfg := config{}
	// for _, opt := range opts {
	// 	opt.apply(&cfg)
	// }
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return func(c *gin.Context) {
		jsonData, _ := ioutil.ReadAll(c.Request.Body)

		// params := c.Params

		logger := logger.GetCurrent().GetLogger("recover")
		c.Set(tracerKey, tracer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		propagator := propagation.TraceContext{}
		parentCtx := propagator.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		logger.Debugf("header extract %v \n", c.Request.Header)
		// logger.Debugf("body %v \n", body)
		// logger.Debugf("params %v \n", params)
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("account", c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := tracer.Start(parentCtx, spanName, opts...)
		span.SetAttributes(attribute.String("body", string(jsonData)))
		propagator.Inject(ctx, propagation.HeaderCarrier(c.Request.Header))
		logger.Debugf("header inject %v \n", c.Request.Header)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		defer func() {

			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")
				if appErr, ok := err.(*common.AppError); ok {
					// trace error
					status := appErr.StatusCode
					checkStatusOfSpan(span, status, c, appErr.Trace)

					appErr.SpanID = span.SpanContext().SpanID().String()
					// handler error
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					panic(err)
					// return
				}

				fmt.Printf("Recover err %v", err)

				appErr := common.ErrInternal(err.(error))

				status := appErr.StatusCode
				checkStatusOfSpan(span, status, c, appErr.Trace)

				appErr.SpanID = span.SpanContext().SpanID().String()
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
				// return
			}
			status := c.Writer.Status()
			checkStatusOfSpan(span, status, c, []string{})
		}()

		c.Next()
	}
}
func checkStatusOfSpan(span oteltrace.Span, status int, c *gin.Context, trace []string) {

	attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
	spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)

	span.SetAttributes(attrs...)
	span.SetStatus(spanStatus, spanMessage)
	if len(trace) > 0 {
		span.SetAttributes(attribute.StringSlice("trace", trace))
	}

	if len(c.Errors) > 0 {
		span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
	}
}
