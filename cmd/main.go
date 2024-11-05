package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/in/http/cv"
	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/in/http/monitoring"
	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/db/cassandra"
	services "github.com/ernstvorsteveld/go-cv-cassandra/domain/serivces"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
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
var (
	g errgroup.Group
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
	apiServer := cv.NewGinCvServer(cv.NewCvApiService(h, c), c)

	monitoringServer := monitoring.NewGinMonitoringServer(monitoring.NewMonitoringApiService(c), c)

	g.Go(func() error {
		return monitoringServer.ListenAndServe()
	})
	slog.Debug("main", "content", "Start Cv API")
	g.Go(func() error {
		return apiServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
