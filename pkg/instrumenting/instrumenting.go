package instrumenting

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpSuccessRegex, _ = regexp.Compile("^2[0-9]{2}$")

	// HTTPRequestsTotalCounter represents a prometheus counter for counting http calls
	HTTPRequestsTotalCounter *kitprometheus.Counter

	// HTTPRequestDurationHistogram represents a promtheus histogram for measuring http calls durations
	HTTPRequestDurationHistogram *kitprometheus.Histogram
)

func init() {
	HTTPRequestsTotalCounter = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests received.",
	}, []string{"component", "handler", "code", "method", "success"})

	HTTPRequestDurationHistogram = kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration in seconds",
	}, []string{"component", "handler", "success"})
}

// ResponseWriter wraps the http.ResponseWriter for adding
// the https response status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader writes the http response status code and capture it inside the writer
func (lrw *ResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// NewResponseWriter implements the ResponseWriter interface and is used
// for capturing the http response status code
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// Middleware increase the counter and record the duration
func Middleware(componentName string, handlerName string, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewResponseWriter(w)

		defer func(begin time.Time) {
			success := httpSuccessRegex.MatchString(strconv.Itoa(lrw.statusCode))
			HTTPRequestsTotalCounter.With("component", componentName, "handler", handlerName, "code", strconv.Itoa(lrw.statusCode), "method", strings.ToLower(r.Method), "success", strconv.FormatBool(success)).Add(1)
			HTTPRequestDurationHistogram.With("component", componentName, "handler", handlerName, "success", strconv.FormatBool(success)).Observe(time.Since(begin).Seconds())
		}(time.Now())

		next.ServeHTTP(lrw, r)
	})
}
