/*
This config package is mainly for loading environment variable (included env file) and json file

Call

	config.Load()

to load the config
*/
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	Postgres string
}

// No need to return error when you can't load the config
// and you are not making a library
func Load() *Config {
	// Add env here
	err := godotenv.Load("./config/.env")

	// Add json config here
	//

	if err != nil {
		panic(err)
	}

	return &Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Postgres: os.Getenv("POSTGRES"),
	}
}
