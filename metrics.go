package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of http requests.",
	},
	// TODO Metrics: Add following labels to the metrics: method, path & app. Uncomment following line and delete the next one.
	// []string{"method", "path", "app"},
	[]string{},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status",
		Help: "Status of http response",
	},
	// TODO Metrics: Add following labels to the metrics: status & app. Uncomment following line and delete the next one.
	// []string{"status", "app"},
	[]string{},
)

var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of http requests.",
	},
	// TODO Metrics: Add following labels to the metrics: method, path & app. Uncomment following line and delete the next one.
	// []string{"method", "path", "app"},
	[]string{},
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx, span := otel.Tracer(appName).Start(r.Context(), "metricsMiddleware")
		defer span.End()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(
			r.Method,
			r.RequestURI,
			appName,
		))
		defer timer.ObserveDuration()
		log := requestLog("metricsMiddleware", r)
		log.Debug("starting metrics...")
		spannedRequest := r.WithContext(spanCtx)
		rw := newResponseWriter(w)
		next.ServeHTTP(w, spannedRequest)
		totalRequests.WithLabelValues(
			r.Method,
			r.RequestURI,
			appName,
		).Inc()
		responseStatus.WithLabelValues(
			strconv.Itoa(rw.statusCode),
			appName,
		).Inc()
		log.Debug("closing metrics...")
	})
}
