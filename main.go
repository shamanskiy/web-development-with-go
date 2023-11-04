package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Shamanskiy/lenslocked/http/server"
	"github.com/Shamanskiy/lenslocked/models"
	"github.com/joho/godotenv"
)

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	server.Run(cfg)
}

func loadEnvConfig() (server.Config, error) {
	var cfg server.Config

	// read .env file and set env variables
	envFile := getEnvFilename()
	err := godotenv.Load(envFile)
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.User == "" || cfg.PSQL.Database == "" {
		return cfg, fmt.Errorf("no postgres config provided")
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func getEnvFilename() string {
	env := os.Getenv("LENSLOCKED_ENV")
	switch env {
	case "PROD":
		return ".env-prod"
	case "":
		return ".env"
	default:
		panic("unknown environment:" + env)
	}
}
