// +build seed

package main

import (
	"log"
	"testing"
)

func TestAlerterSetup(t *testing.T) {
	seedAlerts(t)
}

func seedAlerts(t *testing.T) {
	alerter, err := NewAlerter(*redisUrl, *redisPassword)
	if err != nil {
		log.Fatalf("unable to connect to Redis, %s\n", err.Error())
	}
	alerter.pool.Get().Flush()
	defer alerter.pool.Close()
	result := &SensuResult{
		Client: "localhost",
		Check: SensuCheck{
			Name:         "redis0-redis_role",
			CapsuleName:  "redis0",
			Output:       "expected role master, found role slave",
			Status:       0,
			CapsuleId:    "111111",
			DeploymentId: "987654321",
			AccountSlug:  "compose-test",
		},
	}
	value, err := result.Check.serialize()
	if err != nil {
		log.Fatalf("unable to serialize check, %s\n", err.Error())
	}
	key, field, err := result.Check.alertInfo()
	if err != nil {
		log.Fatalf("unable to get alertInfo, %s\n", err.Error())
	}
	alerter.SetHash(key, field, value)
}
