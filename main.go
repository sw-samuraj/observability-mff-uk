package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	hdrRequestId = "X-Request-ID"
	appAddr      = "0.0.0.0:4040"
	logFile      = "_logs/observability.log"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// Set JSON formatter. Comment out this line to have the text output.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Set logging to a file. Comment out following 2 lines to log on the console.
	f := getLogFile()
	logrus.SetOutput(f)
}

func main() {
	log := funcLog("main")
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.Use(requestLogMiddleware)
	log.Infof("starting observability app on: %s", appAddr)
	http.ListenAndServe(appAddr, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestId(r)
	log := requestLog("homeHandler", requestId)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(hdrRequestId, requestId)
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

func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := requestLog("requestLogMiddleware", getRequestId(r))
		log.Infof("serving request: %s %s%s", r.Method, r.Host, r.RequestURI)
		log.Debugf("user agent: %s", r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func getRequestId(r *http.Request) string {
	requestId := r.Header.Get(hdrRequestId)
	if requestId == "" {
		requestId = uuid.New().String()
		r.Header.Set(hdrRequestId, requestId)
		log := requestLog("getRequestId", requestId)
		log.Warnf("header %s is empty, no request id has been provided", hdrRequestId)
	}
	return requestId
}

func requestLog(f, id string) *logrus.Entry {
	return funcLog(f).WithField("requestId", id)
}

func funcLog(f string) *logrus.Entry {
	return logrus.WithField("func", f)
}

func getLogFile() *os.File {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		funcLog("getLogFile").Fatalf("log file %s can't be created: %v", logFile, err)
	}
	return f
}
