package errorhandling

import (
	"context"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// RecoverFromPanic catches panic from dependencies and recover to an error
func RecoverFromPanic(logger kitlog.Logger, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("recover panic: %v", rec)
				}
				ctx := context.WithValue(context.Background(), kithttp.ContextKeyRequestURI, r.RequestURI)
				ctx = context.WithValue(ctx, kithttp.ContextKeyRequestPath, r.URL.Path)
				ctx = context.WithValue(ctx, kithttp.ContextKeyRequestMethod, r.Method)
				ctx = context.WithValue(ctx, kithttp.ContextKeyRequestUserAgent, r.Header.Get("User-Agent"))
				ctx = context.WithValue(ctx, kithttp.ContextKeyRequestProto, r.Proto)

				kithttp.DefaultErrorEncoder(ctx, Internal("recovered_panic", err), w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
