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
	c *utils.Configuration
)

func getDebugLevel(c *utils.Configuration) slog.Level {
	switch c.DebugLevel {
	case "DEBUG":
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func init() {
	c = &utils.Configuration{}
	c.Read("config", "yml")
	c.Print()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: getDebugLevel(c)}))
	slog.SetDefault(logger)

	slog.Info("Register Handler for Prometheus")
	prometheus.MustRegister(HttpRequestCountWithPath)
}

func main() {

	session := cassandra.NewCassandraSession(c)
	ep := cassandra.NewExperiencePort(c, session)
	tp := cassandra.NewTagPort(c, session)
	h := services.NewCvServices(ep, tp)
	server := cv.NewGinCvServer(cv.NewCvApiService(h, c), c)

	log.Fatal(server.ListenAndServe())
}
