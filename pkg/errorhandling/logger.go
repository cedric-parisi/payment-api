package errorhandling

import (
	"context"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

// Log an error
func Log(ctx context.Context, err error, errLogger kitlog.Logger) {
	var traceID string

	sp := stdopentracing.SpanFromContext(ctx)
	if sp != nil {
		sp.SetTag("error", true)
		sp.SetTag("msg", err.Error())
		defer sp.Finish()

		// If Jaeger is not initialised it's creating panic error
		spCtx := sp.Context()
		switch spCtx.(type) {
		case jaeger.SpanContext:
			traceID = spCtx.(jaeger.SpanContext).TraceID().String()
		}
	}

	// Get HttpStatus code
	statusCode := 0
	if _, ok := err.(kithttp.StatusCoder); ok {
		statusCode = err.(kithttp.StatusCoder).StatusCode()
	}

	errLogger.Log("err", err,
		"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
		"http.path", ctx.Value(kithttp.ContextKeyRequestPath),
		"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
		"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
		"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
		"trace_id", traceID,
		"http.status", statusCode)
}
