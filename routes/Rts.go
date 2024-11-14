package routes

import (
	"net/http"

	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartRoutes(router *mux.Router) {
	StartSemesterRoutes(router)
	StartTimetableRoutes(router)
	StartSubjectClassRoutes(router)
	StartAuthRoutes(router)
	router.HandleFunc("/", controllers.GetIndex).Methods("GET")
	router.Handle("/ping", JWTAuthMiddleware(http.HandlerFunc(Ping))).Methods("GET")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
