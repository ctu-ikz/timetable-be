package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
	"github.com/joho/godotenv"
)

var db *sql.DB

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

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbSSLMode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")
	return db, nil
}

func GetSemesterByTime(time time.Time) (*models.Semester, error) {
	fmt.Println("DB semester requested")

	timeString := time.Format("2006-01-02")

	var semester models.Semester
	err := db.QueryRow(`SELECT id,start,"end",codename FROM "Semester" WHERE $1 BETWEEN "start" AND "end";`, timeString).Scan(&semester.ID, &semester.Start, &semester.End, &semester.Codename)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no semester found")
		}
		return nil, err
	}

	return &semester, nil
}
