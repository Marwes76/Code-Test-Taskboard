package config

import (
	"github.com/joho/godotenv"
	"os"
)

func GetEnvPort() (port string) {
	port = ":" + os.Getenv("PORT")
	return port
}

func GetEnvDbHost() (string) {
	return os.Getenv("DB_HOST")
}

func GetEnvDbPort() (string) {
	return os.Getenv("DB_PORT")
}

func GetEnvDbName() (string) {
	return os.Getenv("DB_NAME")
}

func GetEnvDbUser() (string) {
	return os.Getenv("DB_USER")
}

func GetEnvDbPass() (string) {
	return os.Getenv("DB_PASS")
}

func LoadEnv() (err error) {
	env := os.Getenv("APP_ENV")

	if env == "" || env == "development" {
		// Use development environment files
		err = godotenv.Load(".env.localhost")
	}

	// Production should instead use system environment variables from the running server (for example, a Docker-file), so no
	// .env-file is loaded here

	return err
}