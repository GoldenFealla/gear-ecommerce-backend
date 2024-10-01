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
	defaultHost    = "0.0.0.0"
	defaultPort    = "8080"
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
	godotenv.Load()

	// Add json config here
	//

	hostEnv := os.Getenv("HOST")

	if hostEnv == "" {
		log.Println("env HOST not found, using default Host")
		hostEnv = defaultHost
	}

	portEnv := os.Getenv("PORT")

	if portEnv == "" {
		log.Println("env PORT not found, using default Host")
		portEnv = defaultPort
	}

	posgreSQLEnv := os.Getenv("POSTGRES")

	if posgreSQLEnv == "" {
		log.Println("env POSTGRES not found")
	}

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)

	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}

	timeoutContext := time.Duration(timeout) * time.Second

	return &Config{
		Host:     hostEnv,
		Port:     portEnv,
		Postgres: posgreSQLEnv,
		Timeout:  timeoutContext,
	}
}
