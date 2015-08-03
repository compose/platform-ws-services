package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
 * authHandler takes a handlerfunc, a responsewriter, and an http.Request,
 * checks for auth and called the handler
 */
func authHandler(h HandlerFunc, jw *JsonResponseWriter, r *http.Request) {
	if performAuth(r) {
		h(jw, r)
		jw.t_end = time.Now()
		log.Printf("%s|%s|%s|%db|%d|%dus", r.Header.Get("X-Forwarded-For"), r.Method, r.RequestURI, jw.Size(), jw.Status(), jw.t_end.Sub(jw.t_start)/1000)

	} else {
		jw.WriteHeader(401)
		jw.Write([]byte("401 Unauthorized\n"))
	}
}

/*
 * does this request contain the proper auth header.
 */
func performAuth(r *http.Request) bool {
	if r.Header.Get("Authorization") == "" {
		return false
	}
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) == 2 && s[0] == "Basic" {
		basic_auth, _ := base64.StdEncoding.DecodeString(s[1])
		password := strings.Split(string(basic_auth), ":")[1]
		if string(password) == os.Getenv("COMPOSE_SERVICE_PASSWORD") {
			return true
		}
	}
	return false
}
