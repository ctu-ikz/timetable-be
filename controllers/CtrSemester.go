package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
	"github.com/gorilla/mux"
)

var (
	cache      *models.SemesterCache
	cacheMutex sync.RWMutex
)

func GetSemester(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	cacheMutex.RLock()
	if cache != nil && currentTime.After(cache.Start) && currentTime.Before(cache.End) {
		json.NewEncoder(w).Encode(cache.Semester)
		cacheMutex.RUnlock()
		return
	}
	cacheMutex.RUnlock()

	semester, err := db.GetSemesterByTime(currentTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cacheMutex.Lock()
	cache = &models.SemesterCache{
		Semester: semester,
		Start:    semester.Start,
		End:      semester.End,
	}
	cacheMutex.Unlock()

	json.NewEncoder(w).Encode(semester)
}

func PostSemester(w http.ResponseWriter, r *http.Request) {
	var semester models.Semester
	err := json.NewDecoder(r.Body).Decode(&semester)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newSemester, err := db.PostSemester(&semester)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newSemester)
}

func DeleteSemester(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	if stringid == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(stringid)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	err = db.DeleteSemester(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PutSemester(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	if stringid == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(stringid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var semester models.Semester

	err = json.NewDecoder(r.Body).Decode(&semester)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.PutSemester(id, &semester)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(semester)
}

func GetSemesterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	if stringid == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(stringid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	semester, err := db.GetSemesterByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(semester)
}
