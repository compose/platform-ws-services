package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GroupsHandlerFunc() HandlerFunc {
	return func(w *JsonResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "POST" {
			RegisterGroup(w, r)
		} else {
			UnknownEndpoint(w, r)
		}
	}
}

func RegisterGroup(w *JsonResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var group Group
	err := decoder.Decode(&group)
	if err != nil {
		w.WriteError(err)
		return
	}
	if err = group.Validate(); err != nil {
		w.WriteError(err)
		return
	}
	if err = group.Save(); err != nil {
		w.WriteError(err)
	} else {
		w.WriteOk(201)
	}
}
