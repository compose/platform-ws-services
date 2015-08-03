// +build seed

package main

import (
	"log"
	"testing"

	r "github.com/dancannon/gorethink"
)

func TestSetup(t *testing.T) {
	session, err := initRethinkConn()
	defer session.Close()
	if err != nil {
		log.Fatalf("unable to connect to RethinkDB, %s\n", err.Error())
	}
	setupRethinkdb(session)
	seedRethinkdbChecks(session)
	seedRethinkdbGroups(session)
	seedRethinkdbDeployments(session)
}

func setupRethinkdb(s *r.Session) error {
	_, err := r.DBDrop("alerts").RunWrite(s)
	r.DBCreate("alerts").RunWrite(s)
	return err
}

func seedRethinkdbChecks(s *r.Session) error {
	r.DB("alerts").TableCreate("checks").RunWrite(s)

	data := []*Check{
		&Check{
			Id:          "elasticsearch-es_heap",
			Type:        "elasticsearch",
			Name:        "es_heap",
			Level:       1,
			Title:       "Elasticsearch Heap",
			Description: "total heap usage by an elasticsearch node",
		},
		&Check{
			Id:          "elasticsearch-es_cluster_status",
			Type:        "elasticsearch",
			Name:        "es_heap",
			Level:       1,
			Title:       "Elasticsearch Cluster Status",
			Description: "cluster status (green/yello/red)",
		},
		&Check{
			Id:          "redis-role",
			Type:        "redis",
			Name:        "redis_role",
			Level:       1,
			Title:       "Redis Role",
			Description: "role of redis server (master/slave)",
		},
	}

	s.Use("alerts")
	for _, check := range data {
		r.Table("checks").Insert(check).RunWrite(s)
	}
	r.Table("checks").IndexCreateFunc("type_name", func(row r.Term) interface{} {
		return []interface{}{row.Field("type"), row.Field("name")}
	}).RunWrite(s)
	return nil
}

func seedRethinkdbGroups(s *r.Session) error {
	r.DB("alerts").TableCreate("groups").RunWrite(s)

	r.DB("alerts").Table("groups").Insert(
		&Group{
			Id: "123456789",
			Settings: Settings{
				PagerdutyKey: "123456",
				Slack: Slack{
					ApiKey:  "slack_api_key",
					Channel: "#general",
				},
			},
		},
	).RunWrite(s)
	return nil
}

func seedRethinkdbDeployments(s *r.Session) error {
	r.DB("alerts").TableCreate("deployments").RunWrite(s)

	data := []*Deployment{
		&Deployment{
			Id:      "987654321",
			GroupId: "123456789",
			Type:    "redis",
			Name:    "compose-redis-test",
		},
		&Deployment{
			Id:      "13243647586970",
			GroupId: "123456789",
			Type:    "elasticsearch",
			Name:    "compose-elasticsearch-test",
		},
	}

	for _, deployment := range data {
		r.DB("alerts").Table("deployments").Insert(deployment).RunWrite(s)
	}
	r.DB("alerts").Table("deployments").IndexCreate("group_id").RunWrite(s)
	return nil
}
