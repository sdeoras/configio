package main

import (
	"context"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/configio/simpleconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// create instance of config manager
	// in this case create a local config file manager
	manager := configfile.NewManager(context.Background())

	// create instance of config object
	// config object must satisfy configio.Marshaler interface
	config := new(simpleconfig.Config).Rand()

	// change some config params
	config.PDName = "xyz"

	// marshal out to file
	if err := manager.Marshal(config); err != nil {
		logrus.Fatal(err)
	}
}