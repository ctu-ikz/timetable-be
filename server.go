package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Semester struct {
	ID       int    `json:"id"`
	Year     int    `json:"year"`
	Type     bool   `json:"type"`
	Codename string `json:"codename"`
}

var db *sql.DB

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load connection string components from environment variables
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
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")

	router := mux.NewRouter()
	router.HandleFunc("/semester", getSemester).Methods("GET")
	router.HandleFunc("/dbSemester", getDbSemester).Methods("GET")

	fmt.Println("Server up and running")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getSemester(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Semester requested")
	currentTime := time.Now()
	year := currentTime.Year() % 100

	var subYear int
	if currentTime.Month() > 8 {
		subYear = 1
	} else {
		subYear = 2
	}

	semester := fmt.Sprintf("B%d%d", year, subYear)
	json.NewEncoder(w).Encode(semester)
}

func getDbSemester(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	year := currentTime.Year()

	var subYear int
	if currentTime.Month() > 8 {
		subYear = 1
	} else {
		subYear = 2
	}

	semester, err := dbSemester(db, year, subYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(semester)
}

func dbSemester(db *sql.DB, year int, subyear int) (*Semester, error) {
	fmt.Println("DB semester requested")
	isWinter := subyear == 1
	var semester Semester

	err := db.QueryRow(`SELECT id, year, type, codename FROM "Semester" WHERE year = $1 AND type = $2`, year, isWinter).Scan(&semester.ID, &semester.Year, &semester.Type, &semester.Codename)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no semester found")
		}
		return nil, err
	}
	fmt.Print(semester)

	return &semester, nil
}
