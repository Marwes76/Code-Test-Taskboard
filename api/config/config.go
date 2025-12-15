package config

import (
	"os"
)

func GetEnvApiPort() (port string) {
	port = ":" + os.Getenv("API_PORT")
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
	return os.Getenv("DB_PASSWORD")
}