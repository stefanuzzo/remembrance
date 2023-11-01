package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

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

	t := reflect.TypeOf(c)
	fmt.Println(t)

	nFields := t.NumField()
	fmt.Printf("fields: %d\n", nFields)
	for i := 0; i < nFields; i++ {
		f := t.Field((i))
		fmt.Println(strings.ToLower(f.Name))
	}

	fmt.Printf("mode: %s\n", c.Mode)
}
