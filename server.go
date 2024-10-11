package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/ctu-ikz/timetable-be/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var database *sql.DB

func main() {

	var err error

	database, err = db.ConnectToDB()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/semester", controllers.GetDbSemester).Methods("GET")
	router.HandleFunc("/timetable", controllers.GetThisWeekTimetable).Methods("GET")

	fmt.Println("Server up and running")
	log.Fatal(http.ListenAndServe(":8080", router))
}
