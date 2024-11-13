package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
)

func GetCurrentSubjectClass(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	var semester *models.Semester
	var currentSubjectClass *models.SubjectClass
	classIDQuery := r.URL.Query().Get("class_id")

	if classIDQuery == "" {
		http.Error(w, "Missing class_id parameter", http.StatusBadRequest)
		return
	}

	classID, err := strconv.Atoi(classIDQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	semester, err = db.GetSemesterByTime(currentTime)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if semester == nil {
		http.Error(w, "No semester found", http.StatusInternalServerError)
		return
	}

	currentSubjectClass, err = db.GetCurrentSubjectClass(*semester, classID, currentTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No subject class found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(currentSubjectClass)
}
