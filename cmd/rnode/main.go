package main

import (
	"fmt"
	"os"

	"github.com/stefanuzzo/internal/configuration"
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
}
