package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	hdrRequestId = "X-Request-ID"
	appAddr      = "0.0.0.0:4040"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
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
