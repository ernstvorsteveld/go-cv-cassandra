package main

import (
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/adapter/in/http/cv"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/adapter/out/db/cassandra"
	services "github.com/ernstvorsteveld/go-cv-cassandra/src/domain/serivces"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
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

	// petStore := api.NewPetStore()
	// s := NewGinPetServer(petStore, *port)
	// // And we serve HTTP until the world ends.
	log.Fatal(server.ListenAndServe())
}
