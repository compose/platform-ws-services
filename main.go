package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	r "github.com/dancannon/gorethink"
)

var (
	listenPort      = flag.String("p", "8000", "api listen port")
	redisUrl        = flag.String("redisUrl", os.Getenv("REDIS_URL"), "Redis URL")
	redisPassword   = flag.String("redisPassword", os.Getenv("REDIS_PASSWORD"), "Redis password")
	rethinkUrl      = flag.String("rethinkUrl", os.Getenv("RETHINK_URL"), "RethinkDB URL")
	rethinkPassword = flag.String("rethinkPassword", os.Getenv("RETHINK_PASSWORD"), "RethinkDB password")
	session         *r.Session
)

func main() {
	flag.Parse()

	log.SetFlags(0)

	alerter, err := NewAlerter(*redisUrl, *redisPassword)
	failOnError(err, fmt.Sprintf("Failed to connect to Redis (%s)", *redisUrl))

	session, err = initRethinkConn()
	failOnError(err, fmt.Sprintf("Failed to connect to RethinkdB (%s)", *rethinkUrl))
	defer session.Close()
	session.Use("alerts")

	notifier := NewNotifier()
	notifier.Start()

	listener, err := NewListener(alerter)
	failOnError(err, "Failed to create a listener")
	listener.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("/status", ApiStatusHandler)
	mux.Handle("/margo/deployments", http.StripPrefix("/margo/deployments", DeploymentsHandlerFunc(alerter)))
	mux.Handle("/margo/deployments/", http.StripPrefix("/margo/deployments/", DeploymentsHandlerFunc(alerter)))
	mux.Handle("/margo/groups", http.StripPrefix("/margo/groups", GroupsHandlerFunc()))
	mux.Handle("/margo/checks", http.StripPrefix("/margo/checks", ChecksHandlerFunc()))
	mux.Handle("/margo/checks/", http.StripPrefix("/margo/checks/", ChecksHandlerFunc()))
	mux.Handle("/margo/alerts", http.StripPrefix("/margo/alerts", AlertsHandlerFunc(listener)))

	log.Printf("listening on %s\n", *listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *listenPort), mux))
}

func initRethinkConn() (*r.Session, error) {
	addresses := strings.Split(*rethinkUrl, ",")
	connectOpts := r.ConnectOpts{
		Addresses:     addresses,
		Timeout:       time.Duration(10 * time.Second),
		MaxIdle:       10,
		MaxOpen:       10,
		DiscoverHosts: false,
	}
	if *rethinkPassword != "" {
		connectOpts.AuthKey = *rethinkPassword
	}
	return r.Connect(connectOpts)
}

func ApiStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func UnknownEndpoint(w *JsonResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
