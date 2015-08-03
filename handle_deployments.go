package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
)

// TOOD: refactor
func DeploymentsHandlerFunc(alerter *Alerter) HandlerFunc {
	return func(w *JsonResponseWriter, r *http.Request) {
		lookup := html.EscapeString(r.URL.Path)
		urlParts := strings.Split(lookup, "/")
		if urlParts[0] == "" {
			if strings.ToUpper(r.Method) == "POST" {
				RegisterDeployment(w, r)
			} else {
				UnknownEndpoint(w, r)
			}
		} else if len(urlParts) == 1 && strings.ToUpper(r.Method) == "GET" && urlParts[0] != "" {
			GetDeployment(w, r, urlParts[0], alerter)
		} else {
			deployment, err := LookupDeploymentById(urlParts[0])
			if err != nil {
				w.WriteError(errors.New("deployment not found"))
				return
			}
			if len(urlParts) == 1 {
				if strings.ToUpper(r.Method) == "DELETE" {
					UnregisterDeployment(w, r, deployment)
				} else {
					UnknownEndpoint(w, r)
				}
			} else if strings.Contains(lookup, "/checks") {
				if strings.ToUpper(r.Method) == "GET" {
					if len(urlParts) == 2 {
						GetDeploymentChecks(w, r, deployment)
					} else {
						GetDeploymentCheck(w, r, deployment, urlParts[2])
					}
				} else {
					UnknownEndpoint(w, r)
				}
			} else {
				UnknownEndpoint(w, r)
			}
		}
	}
}

func GetDeployment(w *JsonResponseWriter, r *http.Request, deployment string, alerter *Alerter) {
	alerts, err := alerter.GetAll(fmt.Sprintf("%s:*", deployment))
	if err != nil {
		w.WriteError(err)
	} else {
		w.WriteJson(alerts)
	}
}

func RegisterDeployment(w *JsonResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var deployment Deployment
	err := decoder.Decode(&deployment)
	if err != nil {
		w.WriteError(err)
		return
	}
	if err = deployment.Validate(); err != nil {
		w.WriteError(err)
		return
	}
	if err = deployment.Save(); err != nil {
		w.WriteError(err)
	} else {
		w.WriteOk(201)
	}
}

func UnregisterDeployment(w *JsonResponseWriter, r *http.Request, deployment Deployment) {
	if err := deployment.Delete(); err != nil {
		w.WriteError(err)
	} else {
		w.WriteOk(200)
	}
}

func GetDeploymentChecks(w *JsonResponseWriter, r *http.Request, deployment Deployment) {
	checks, err := deployment.CurrentChecks()
	if err != nil {
		w.WriteError(err)
	} else {
		w.WriteJson(checks)
	}
}

func GetDeploymentCheck(w *JsonResponseWriter, r *http.Request, deployment Deployment, name string) {
	check, err := deployment.CheckByName(name)
	if err != nil {
		w.WriteError(err)
	} else {
		w.WriteJson(check)
	}
}
