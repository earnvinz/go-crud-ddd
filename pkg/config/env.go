package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

func LoadEnv() DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return DBConfig{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
	}
}

func ConnectDB(cfg DBConfig) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	return db
}
