package utils

import (
	"os"

	log "github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	DB DBConfiguration
}

type DBConfiguration struct {
	Cassandra CassandraConfiguration
}

type CassandraConfiguration struct {
	Keyspace string
	Url      string
}

type ConfigurationManager interface {
	Read(fname string, ftype string)
	Print()
}

func (c *Configuration) Read(fname string, ftype string) {
	c.readFile(fname, ftype)
	c.readEnvironment()
}

const CASSANDRA_URL = "CASSANDRA_URL"
const CASSANDRA_KEYSPACE = "CASSANDRA_KEYSPACE"

func (c *Configuration) readFile(fname string, ftype string) {
	viper.SetConfigName(fname)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType(ftype)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
	}
}

func (c *Configuration) readEnvironment() {
	value := os.Getenv(CASSANDRA_URL)
	if value != "" {
		log.Infof("Url environment used: %s.", value)
		c.DB.Cassandra.Url = value
	}

	value = os.Getenv(CASSANDRA_KEYSPACE)
	if value != "" {
		log.Infof("Keyspace environment used: %s.", value)
		c.DB.Cassandra.Keyspace = value
	}
}

func (c *Configuration) Print() {
	log.Info("Configuration:\n")
	log.Infof("%+v\n", *c)
}
