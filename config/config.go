/*
This config package is mainly for loading environment variable (included env file) and json file

Call

	config.Load()

to load the config
*/
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 30
)

type Config struct {
	Host     string
	Port     string
	Postgres string
	Timeout  time.Duration
}

// No need to return error when you can't load the config
// and you are not making a library
func Load() *Config {
	// Add env here
	err := godotenv.Load()

	// Add json config here
	//

	if err != nil {
		panic(err)
	}

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)

	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}

	timeoutContext := time.Duration(timeout) * time.Second

	return &Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Postgres: os.Getenv("POSTGRES"),
		Timeout:  timeoutContext,
	}
}
