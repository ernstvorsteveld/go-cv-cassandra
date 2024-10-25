package main

import (
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/adapter/in/http/cv"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/adapter/out/db/cassandra"
	services "github.com/ernstvorsteveld/go-cv-cassandra/src/domain/serivces"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_with_path",
			Help: "Number of HTTP requests by path.",
		},
		[]string{"url"},
	)
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	log.Debugf("Register Handler for Prometheus")
	prometheus.MustRegister(HttpRequestCountWithPath)
}

func main() {
	c := &utils.Configuration{}
	c.Read("config", "yml")
	c.Print()

	session := cassandra.NewCassandraSession(c)
	ep := cassandra.NewExperiencePort(c, session)
	tp := cassandra.NewTagPort(c, session)
	h := services.NewCvServices(ep, tp)
	server := cv.NewGinCvServer(cv.NewCvApiService(h), c.Api.Port)

	log.Fatal(server.ListenAndServe())
}
