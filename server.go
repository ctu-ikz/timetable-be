package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Semester struct {
	ID       int    `json:"id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Codename string `json:"codename"`
}

var database *sql.DB

func main() {

	var err error

	database, err = db.ConnectToDB()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()	

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

	semester, err := dbSemester(database, currentTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(semester)
}

func dbSemester(db *sql.DB, time time.Time) (*Semester, error) {
	fmt.Println("DB semester requested")

	timeString := time.Format("2006-01-02")

	var semester Semester
	err := db.QueryRow(`SELECT id,start,"end",codename FROM "Semester" WHERE $1 BETWEEN "start" AND "end";`, timeString).Scan(&semester.ID, &semester.Start, &semester.End, &semester.Codename)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no semester found")
		}
		return nil, err
	}

	return &semester, nil
}