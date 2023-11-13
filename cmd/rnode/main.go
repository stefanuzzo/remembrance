package main

import (
	"fmt"
	"os"

	"github.com/stefanuzzo/internal/cluster"
	"github.com/stefanuzzo/internal/configuration"
	"github.com/stefanuzzo/internal/key"
)

const configurationFilePath = "../../config/config.json"

func main() {
	fmt.Println("hello, world")

	c, err := configuration.GetConfiguration(false, configurationFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("mode: %s\n", c.Mode)

	err = key.InitializeKeys(c.KeysDirectory)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	uuid, generated, err := cluster.GetOrSetNodeId(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	var sGenerated string
	if generated {
		sGenerated = "new"
	} else {
		sGenerated = "pre-generated"
	}

	fmt.Printf("Node id: %s (%s)\n", uuid.String(), sGenerated)
}
