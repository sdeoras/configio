package main

import (
	"context"

	"fmt"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/configio/simpleconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// create instance of config manager
	// in this case create a local config file manager
	manager, err := configfile.NewManager(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	// create instance of config object
	// config object must satisfy configio.Marshaler interface
	config := new(simpleconfig.Config).Rand()

	// change some config params
	config.Name = "xyz"
	config.Value = 7
	config.ReadOnly = true

	if jb, err := config.Marshal(); err != nil {
		logrus.Fatal(err)
	} else {
		fmt.Println("config struct:")
		fmt.Println(string(jb))
	}

	fmt.Println("saving to backend")
	// marshal out to file
	if err := manager.Marshal(config); err != nil {
		logrus.Fatal(err)
	}

	if err := manager.Unmarshal(config); err != nil {
		logrus.Fatal(err)
	}

	if jb, err := config.Marshal(); err != nil {
		logrus.Fatal(err)
	} else {
		fmt.Println("reading back:")
		fmt.Println(string(jb))
	}
}
