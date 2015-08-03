package main

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"
)

func ChecksHandlerFunc() HandlerFunc {
	return func(w *JsonResponseWriter, r *http.Request) {
		lookup := html.EscapeString(r.URL.Path)
		if strings.ToUpper(r.Method) == "POST" {
			RegisterCheck(w, r)
		} else if strings.ToUpper(r.Method) == "DELETE" && lookup != "" {
			DeleteCheck(w, r, lookup)
		} else if strings.ToUpper(r.Method) == "GET" && lookup != "" {
			GetChecks(w, r, lookup)
		} else {
			UnknownEndpoint(w, r)
		}
	}
}

func RegisterCheck(w *JsonResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var check Check
	err := decoder.Decode(&check)
	if err != nil {
		w.WriteError(err)
		return
	}
	if err = check.Validate(); err != nil {
		w.WriteError(err)
		return
	}
	if err = check.Save(); err != nil {
		w.WriteError(err)
	} else {
		w.WriteOk(201)
	}
}

func DeleteCheck(w *JsonResponseWriter, r *http.Request, check_type string) {
	toDelete := &Check{
		Id: check_type,
	}
	if err := toDelete.Delete(); err != nil {
		w.WriteError(err)
	} else {
		w.WriteOk(200)
	}
}

func GetChecks(w *JsonResponseWriter, r *http.Request, check_type string) {
	checks, err := GetChecksByType(check_type)
	if err != nil {
		w.WriteError(err)
	} else {
		w.WriteJson(checks)
	}
}
