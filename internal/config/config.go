package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPServer
	Postgres
}

type HTTPServer struct {
	ServerAddress string
	Timeout       time.Duration
	IdleTimeout   time.Duration
}

type Postgres struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

func MustLoadConfig() *Config {
	pgDb := os.Getenv("POSTGRES_DB")
	if pgDb == "" {
		log.Fatal("POSTGRES_DB environment variable not set")
	}

	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		log.Fatal("POSTGRES_HOST environment variable not set")
	}

	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		log.Fatal("POSTGRES_PORT environment variable not set")
	}

	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		log.Fatal("POSTGRES_USER environment variable not set")
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		log.Fatal("POSTGRES_PASSWORD environment variable not set")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	serverAddress = serverAddress + ":" + appPort

	timeoutStr := os.Getenv("SERVER_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "4s"
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalf("invalid SERVER_TIMEOUT: %v", err)
	}

	idleTimeoutStr := os.Getenv("SERVER_IDLE_TIMEOUT")
	if idleTimeoutStr == "" {
		idleTimeoutStr = "60s"
	}
	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		log.Fatalf("invalid SERVER_IDLE_TIMEOUT: %v", err)
	}

	srv := HTTPServer{
		ServerAddress: serverAddress,
		Timeout:       timeout,
		IdleTimeout:   idleTimeout,
	}

	pg := Postgres{
		Host:     pgHost,
		Port:     pgPort,
		DB:       pgDb,
		User:     pgUser,
		Password: pgPassword,
	}

	cfg := Config{
		HTTPServer: srv,
		Postgres:   pg,
	}

	return &cfg
}
