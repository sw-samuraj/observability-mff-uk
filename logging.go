package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx, span := otel.Tracer(appName).Start(r.Context(), "loggingMiddleware")
		defer span.End()
		log := requestLog("loggingMiddleware", r)
		log.Infof("serving request: %s %s%s", r.Method, r.Host, r.RequestURI)
		log.Debugf("user agent: %s", r.UserAgent())
		spannedRequest := r.WithContext(spanCtx)
		w.Header().Set(hdrRequestId, getRequestId(r))
		w.Header().Set(hdrCorrelationId, getCorrelationId(r))
		next.ServeHTTP(w, spannedRequest)
	})
}

func requestLog(f string, r *http.Request) *logrus.Entry {
	// TODO Logging: Obtain requestId, correlationId & traceId to the log entry. Uncomment following 3 lines.
	// rid := getRequestId(r)
	// cid := getCorrelationId(r)
	// tid := getTracingId(r)
	return funcLog(f).WithFields(logrus.Fields{
		// TODO Logging: Add requestId, correlationId & traceId to the log entry. Uncomment following 3 lines.
		// "requestId":     rid,
		// "correlationId": cid,
		// "traceId":       tid,
	})
}

func getRequestId(r *http.Request) string {
	requestId := r.Header.Get(hdrRequestId)
	if requestId == "" {
		requestId = uuid.New().String()
		r.Header.Set(hdrRequestId, requestId)
		log := requestLog("getRequestId", r)
		log.Warnf("header %s is empty, no request id has been provided", hdrRequestId)
	}
	return requestId
}

func getCorrelationId(r *http.Request) string {
	correlationId := r.Header.Get(hdrCorrelationId)
	if correlationId == "" {
		correlationId = uuid.New().String()
		r.Header.Set(hdrCorrelationId, correlationId)
		log := requestLog("getCorrelationId", r)
		log.Warnf("header %s is empty, no correlation id has been provided", hdrCorrelationId)
	}
	return correlationId
}

func getTracingId(r *http.Request) string {
	return r.Header.Get(hdrTracingId)
}

func funcLog(f string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":  appName,
		"func": f,
		// TODO Vita: add hostname and/or ip address
	})
}

func getLogFile() *os.File {
	file := logFile
	if appName != defaultAppName {
		file = fmt.Sprintf(logFilePattern, appName)
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		funcLog("getLogFile").Fatalf("log file %s can't be created: %v", logFile, err)
	}
	return f
}
