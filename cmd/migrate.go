package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var migrationsDir = "../migrations"

func checkAndCreateDB(dbUser, dbPass, dbHost, dbPort, dbName string) error {
	adminDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", dbUser, dbPass, dbHost, dbPort)
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to admin DB: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping admin DB: %w", err)
	}

	var exists bool
	err = adminDB.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		fmt.Printf("Database %s does not exist. Creating...\n", dbName)
		_, err := adminDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Println("Database created successfully.")
	} else {
		fmt.Println("Database already exists.")
	}

	return nil
}

func connectAppDB(dbUser, dbPass, dbHost, dbPort, dbName string) (*sql.DB, error) {
	appDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	appDB, err := sql.Open("postgres", appDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to app DB: %w", err)
	}
	if err := appDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping app DB: %w", err)
	}
	fmt.Println("Connected to application database successfully!")
	return appDB, nil
}

func runMigration(dbUser, dbPass, dbHost, dbPort, dbName string) error {
	appDB, err := connectAppDB(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return err
	}
	defer appDB.Close()

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsDir, file.Name()))
		}
	}
	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		fmt.Println("Running migration:", file)
		sqlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}
		if _, err := appDB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	fmt.Println("All migrations ran successfully.")
	return nil
}

func rollbackMigration(dbUser, dbPass, dbHost, dbPort, dbName string) error {
	appDB, err := connectAppDB(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return err
	}
	defer appDB.Close()

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".down.sql") {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsDir, file.Name()))
		}
	}

	// Rollback ต้อง reverse order
	sort.Sort(sort.Reverse(sort.StringSlice(sqlFiles)))

	for _, file := range sqlFiles {
		fmt.Println("Rolling back:", file)
		sqlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}
		if _, err := appDB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute rollback %s: %w", file, err)
		}
	}

	fmt.Println("Rollback completed successfully.")
	return nil
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	if err := checkAndCreateDB(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
		log.Fatal(err)
	}

	// รับ argument เช่น "up" หรือ "down"
	if len(os.Args) < 2 {
		log.Fatal("Missing action. Use: go run main.go up OR down")
	}

	action := os.Args[1]

	switch action {
	case "up":
		if err := runMigration(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := rollbackMigration(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown action: %s. Use: up or down", action)
	}
}
