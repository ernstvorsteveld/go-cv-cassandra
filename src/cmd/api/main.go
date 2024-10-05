package main

import "github.com/ernstvorsteveld/go-cv-cassandra/src/utils"

func main() {
	c := utils.Configuration{}
	c.Read("config", "yml")
	c.Print()
}
