package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/ctu-ikz/timetable-be/db"
	"github.com/gorilla/mux"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer db.GetDB().Close()

	router := mux.NewRouter()
	controllers.StartRoutes(router)

	fmt.Println("Server up and running")
	log.Fatal(http.ListenAndServe(":8080", router))
}
