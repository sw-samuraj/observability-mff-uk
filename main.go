package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	// "log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	defaultAppName    = "my-app"
	defaultAppPort    = "4040"
	defaultTracingUrl = "http://localhost:14268/api/traces"
	hdrCorrelationId  = "X-Correlation-ID"
	hdrRequestId      = "X-Request-ID"
	hdrTracingId      = "X-Tracing-ID"
	hdrUserAgent      = "User-Agent"
	logFile           = "_logs/observability.log"
	logFilePattern    = "_logs/%s.log"
	// tracingUrl     = "http://grafana.edu.dobias.info:14268/api/traces"
)

var (
	appAddr       string
	appName       string
	appPort       string
	downstreamUrl string
	tracingUrl    string
	goVersion     = runtime.Version()
)

func init() {
	// parse command line arguments
	flag.StringVar(&appName, "n", defaultAppName, "Application name.")
	flag.StringVar(&appPort, "p", defaultAppPort, "Application port.")
	flag.StringVar(&downstreamUrl, "d", "", "Downstream URL. Empty string triggers no call to downstream service.")
	flag.StringVar(&tracingUrl, "t", defaultTracingUrl, "Tracing URL.")
	printHelp := flag.Bool("h", false, "Print help.")
	flag.Parse()
	if *printHelp {
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	appAddr = fmt.Sprintf("0.0.0.0:%s", appPort)
	// Set logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// TODO Logging: Enable debug level of logging. Uncomment following line.
	// logrus.SetLevel(logrus.DebugLevel)
	// Set JSON formatter.
	// TODO Logging: Enable json formatting for logging. Uncomment following line.
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// Set logging to a file. Comment out following 2 lines to log on the console.
	// TODO Logging: Enable logging to a file. Uncomment following 2 lines.
	// f := getLogFile()
	// logrus.SetOutput(f)
	// Register Prometheus metrics
	// TODO Metrics: Register metrics. Uncomment following 3 lines.
	// prometheus.Register(totalRequests)
	// prometheus.Register(responseStatus)
	// prometheus.Register(httpDuration)
	// Set tracing provider
	// TODO Tracing: Configure tracing provider. Uncomment following 4 lines.
	// tp, err := tracerProvider(tracingUrl)
	// if err != nil {
	//	log.Fatal(err)
	// }
	// TODO Tracing: Enable tracing provider. Uncomment following line.
	// otel.SetTracerProvider(tp)
	// TODO Tracing: Enable cross-boundary context propagation for tracing. Uncomment following line.
	// otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	log := funcLog("main")
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	// TODO Metrics: Expose metrics endpoint. Uncomment following line.
	// r.Handle("/metrics", promhttp.Handler())
	// TODO Tracing: Enable tracing middleware. Uncomment following line.
	// r.Use(tracingMiddleware)
	// TODO Metrics: Enable metrics middleware. Uncomment following line.
	// r.Use(metricsMiddleware)
	// TODO Logging: Enable logging middleware. Uncomment following line.
	// r.Use(loggingMiddleware)
	log.Infof("starting observability app on: %s", appAddr)
	http.ListenAndServe(appAddr, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(appName).Start(r.Context(), "homeHandler")
	defer span.End()
	log := requestLog("homeHandler", r)
	if downstreamUrl != "" {
		callDownstream(r, spanCtx, downstreamUrl)
	}
	randomizeLatency()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	resp["message"] = "Observability check: 👌"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("can't marshal json: %v", err)
	}
	log.Infof("writing response with status: %d", http.StatusOK)
	w.Write(jsonResp)
	return
}

func randomizeLatency() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func callDownstream(r *http.Request, ctx context.Context, url string) {
	spanCtx, span := otel.Tracer(appName).Start(ctx, "callDownstream")
	defer span.End()
	log := requestLog("callDownstream", r)
	log.Infof("calling downstream service: %s", url)
	downstreamRequest := getDownstreamRequest(spanCtx, r, url)
	client := getClient()
	resp, err := client.Do(downstreamRequest)
	if err != nil {
		log.Errorf("error calling downstream service: %s", err)
	}
	log.Infof("downstream service returned http code: %d", resp.StatusCode)
	log.Infof("downstream service returned request id: %s", resp.Header.Get(hdrRequestId))
	log.Debugf("downstream service returned correlation id: %s", resp.Header.Get(hdrCorrelationId))
}

func getDownstreamRequest(ctx context.Context, r *http.Request, url string) *http.Request {
	downstreamRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(downstreamRequest.Header))
	if err != nil {
		log := requestLog("getDownstreamRequest", r)
		log.Errorf("error assembling downstream request: %s", err)
	}
	downstreamRequest.Header.Set(hdrUserAgent, fmt.Sprintf("Golang/%s", goVersion))
	downstreamRequest.Header.Add(hdrRequestId, uuid.New().String())
	downstreamRequest.Header.Add(hdrCorrelationId, getCorrelationId(r))
	return downstreamRequest
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
