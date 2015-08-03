package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	Resolved = iota // 0
	Warning         // 1
	Critical        // 2
	Unknown         // 3
)

type Listener struct {
	SensuChan chan *SensuResult
	alerter   *Alerter
}

type SensuResult struct {
	Client string     `json:"client"` // the client the check came from
	Check  SensuCheck `json:"check"`  // the check info
}

type SensuCheck struct {
	Name         string  `json:"name"`
	CapsuleName  string  `json:"capsule_name"`
	Output       string  `json:"output"`
	Status       float64 `json:"status"`
	CapsuleId    string  `json:"capsule_id",omitempty"`
	DeploymentId string  `json:"deployment_id,omitempty"`
	AccountSlug  string  `json:"account,omitempty"`
}

func (r SensuResult) String() string {
	return fmt.Sprintf("%s - %s", r.Client, r.Check.Name)
}

func NewListener(alerter *Alerter) (*Listener, error) {
	return &Listener{
		SensuChan: make(chan *SensuResult),
		alerter:   alerter,
	}, nil
}

func (l *Listener) Start() {
	go l.listenForChecks()
}

func (l *Listener) listenForChecks() {
	for {
		select {
		case result := <-l.SensuChan:
			if err := l.processSensuResult(result); err != nil {
				log.Printf("ERROR: Unable to process sensu result, %s", err.Error())
			}
		case <-time.After(100 * time.Millisecond):
			// NOP, just breath
		}
	}
}

func (l *Listener) processSensuResult(result *SensuResult) error {
	if result.Check.Status != Resolved {
		value, _ := result.Check.serialize()
		key, field, err := result.Check.alertInfo()
		if err != nil {
			log.Printf("Unable to act on the following check:\n%s\n", result)
		} else {
			l.alerter.SetHash(key, field, value)
		}
	} else {
		key, field, err := result.Check.alertInfo()
		if err == nil {
			l.alerter.Resolve(key, field)
		}
	}
	return nil
}

func (c *SensuCheck) alertInfo() (string, string, error) {
	if c.CapsuleId == "" {
		return "", "", errors.New("capsule_id required")
	} else if c.DeploymentId == "" {
		return "", "", errors.New("deployment_id required")
	}
	latterPart := strings.TrimPrefix(c.Name, c.CapsuleName)
	checkName := strings.Replace(latterPart, "-", "", 1)
	key := fmt.Sprintf("%s:%s", c.DeploymentId, c.CapsuleId)
	return key, checkName, nil
}

func (c *SensuCheck) serialize() ([]byte, error) {
	return json.Marshal(c)
}
