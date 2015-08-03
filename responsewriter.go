package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

/*
 * responsewriter that has some convenience methods, and tracks the status and body size
 * for later consumption by a logger
 */
type JsonResponseWriter struct {
	w       http.ResponseWriter
	status  int
	size    int
	t_start time.Time
	t_end   time.Time
}

func (j *JsonResponseWriter) WriteError(err error) {
	log.Printf("%s", err.Error())
	j.WriteHeader(http.StatusBadRequest)
	j.WriteJson(map[string]string{"error": err.Error()})
}

func (j *JsonResponseWriter) WriteJson(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		j.WriteError(err)
		return
	}
	j.Write(data)
}

func (j *JsonResponseWriter) WriteOk(s int) {
	if s == 0 {
		j.WriteHeader(http.StatusOK)
	} else {
		j.WriteHeader(s)
	}
	j.WriteJson(map[string]int{"ok": 1})
}

func (j *JsonResponseWriter) Header() http.Header {
	return j.w.Header()
}

func (j *JsonResponseWriter) Write(b []byte) (int, error) {
	if j.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		j.status = http.StatusOK
	}
	size, err := j.w.Write(b)
	j.size += size

	return size, err
}

func (j *JsonResponseWriter) WriteHeader(s int) {
	j.w.WriteHeader(s)
	j.status = s
}

func (j *JsonResponseWriter) Status() int {
	return j.status
}

func (j *JsonResponseWriter) Size() int {
	return j.size
}

/*
 * generic handler that uses the JsonResponseWriter,
 * and ensures that we have auth happening
 */
type HandlerFunc func(*JsonResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jw := &JsonResponseWriter{w: w}
	jw.t_start = time.Now()
	jw.Header().Set("Content-Type", "application/json;charset=utf-8")

	authHandler(h, jw, r)
}
