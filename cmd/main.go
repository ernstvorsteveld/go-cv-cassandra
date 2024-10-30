package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/in/http/cv"
	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/db/cassandra"
	services "github.com/ernstvorsteveld/go-cv-cassandra/domain/serivces"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
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
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	slog.Info("Register Handler for Prometheus")
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
	server := cv.NewGinCvServer(cv.NewCvApiService(h, c), c)

	log.Fatal(server.ListenAndServe())
}
