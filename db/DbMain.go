package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() error {
	var err error
	db, err = ConnectToDB()
	if err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}

	err = applyMigrations()
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}

	return nil
}

func applyMigrations() error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create db instance: %v", err)
	}

	migrationsDir := "/migrations"
	fmt.Println("Migrations directory: ", migrationsDir)

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %v", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)

	if err != nil {
		return fmt.Errorf("could not start migration: %v", err)
	}

	if err := migrateInstance.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("an error occurred during migration: %v", err)
	}

	fmt.Println("Migrations successfully applied")
	return nil
}

func ConnectToDB() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	fmt.Println("Connecting to DB")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping db: %v", err)
	}

	fmt.Println("Database connected on " + dbHost)
	return db, nil
}

func GetDB() *sql.DB {
	return db
}
