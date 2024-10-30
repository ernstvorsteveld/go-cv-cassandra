package utils

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	DB  DBConfiguration
	Api APIConfiguration
}

type DBConfiguration struct {
	Cassandra CassandraConfiguration
}

type CassandraConfiguration struct {
	Keyspace string
	Url      string
	Port     string
	Retries  int8
	Username string
	Secret   SensitiveInfo
}

type APIConfiguration struct {
	Port          string
	Expectedhosts []string
}

type SensitiveInfo string

func (s SensitiveInfo) String() string {
	return "****"
}

func (s SensitiveInfo) Value() string {
	return string(s)
}

type ConfigurationManager interface {
	Read(fname string, ftype string)
	Print()
}

func (c *Configuration) Read(fname string, ftype string) {
	c.readFile(fname, ftype)
	c.readEnvironment()

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(c)
	if err != nil {
		fmt.Printf("error while reading configuration %v", err)
		panic(err)
	}
}

const CASSANDRA_URL = "CASSANDRA_URL"
const CASSANDRA_KEYSPACE = "CASSANDRA_KEYSPACE"
const CASSANDRA_SECRET = "CASSANDRA_SECRET"

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

	value = os.Getenv(CASSANDRA_SECRET)
	if value != "" {
		log.Infof("Secret environment used: *********.")
		c.DB.Cassandra.Secret = SensitiveInfo(value)
	}
}

func (c *Configuration) Print() {
	log.Infof("Configuration:%+v", *c)
}
