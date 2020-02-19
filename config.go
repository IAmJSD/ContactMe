package main

import (
	"os"
	"time"

	"github.com/jakemakesstuff/structuredhttp"
)

// Config defines the main configuration.
var Config BaseStructure

// Initialises the config.
func init() {
	ConfigURL := os.Getenv("CONFIG_URL")
	if ConfigURL == "" {
		panic("Config URL not found")
	}
	r, err := structuredhttp.GET(ConfigURL).Timeout(time.Second * 10).Run()
	if err != nil {
		panic(err)
	}
	j, err := r.JSON()
	if err != nil {
		panic(err)
	}
	Config = j.(BaseStructure)
}
