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
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultHost    = "0.0.0.0"
	defaultPort    = "8080"
	defaultTimeout = 30
)

type S3Config struct {
	AccountID        string
	AccountKeyID     string
	AccountKeySecret string
}

type Config struct {
	Host         string
	Port         string
	Postgres     string
	Redis        string
	Timeout      time.Duration
	S3           *S3Config
	AllowOrigins []string
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
		log.Fatalln("env POSTGRES not found")
	}

	RedisEnv := os.Getenv("REDIS")

	if RedisEnv == "" {
		log.Fatalln("env REDIS not found")
	}

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)

	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}

	timeoutContext := time.Duration(timeout) * time.Second

	allowOriginsStr := os.Getenv("ALLOW_ORIGINS")
	allowOrigins := strings.Split(allowOriginsStr, ";")

	if allowOriginsStr != "" {
		log.Println("received ALLOW_ORIGINS env:")

		for _, v := range allowOrigins {
			log.Println(v)
		}

	}

	S3AccountID := os.Getenv("S3_ACCOUNT_ID")
	S3AccountKeyID := os.Getenv("S3_ACCOUNT_KEY_ID")
	S3AccountKeySecret := os.Getenv("S3_ACCOUNT_KEY_SECRET")

	if S3AccountID == "" {
		log.Fatalln("S3_ACCOUNT_ID not found, Please add one")
	}

	if S3AccountKeyID == "" {
		log.Fatalln("S3_ACCOUNT_KEY_ID not found, Please add one")
	}

	if S3AccountKeySecret == "" {
		log.Fatalln("S3_ACCOUNT_KEY_SECRET not found, Please add one")
	}

	return &Config{
		Host:     hostEnv,
		Port:     portEnv,
		Postgres: posgreSQLEnv,
		Redis:    RedisEnv,
		Timeout:  timeoutContext,
		S3: &S3Config{
			AccountID:        S3AccountID,
			AccountKeyID:     S3AccountKeyID,
			AccountKeySecret: S3AccountKeySecret,
		},
		AllowOrigins: allowOrigins,
	}
}
