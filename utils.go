package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file")
		panic(err)
	}

	envVar := os.Getenv(key)
	if envVar == "" {
		log.Fatalf("Unable to load environment variable %s", key)
	}

	return envVar
}
