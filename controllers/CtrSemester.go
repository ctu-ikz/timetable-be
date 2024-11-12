package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
	"github.com/gorilla/mux"
)

func GetSemester(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	SemesterCache.Mutex.RLock()
	for _, semester := range SemesterCache.Data {
		if currentTime.After(semester.Start) && currentTime.Before(semester.End) {
			json.NewEncoder(w).Encode(semester)
			SemesterCache.Mutex.RUnlock()
			return
		}
	}
	SemesterCache.Mutex.RUnlock()

	semester, err := db.GetSemesterByTime(currentTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SemesterCache.Mutex.Lock()
	if semester.ID != nil {
		SemesterCache.Data[*semester.ID] = *semester
	} else {
		http.Error(w, "Semester ID is nil", http.StatusInternalServerError)
	}
	SemesterCache.Mutex.Unlock()

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

	SemesterCache.Mutex.Lock()
	if newSemester.ID != nil {
		SemesterCache.Data[*newSemester.ID] = *newSemester
		SemesterCache.Mutex.Unlock()
	} else {
		http.Error(w, "New Semester ID is nil", http.StatusInternalServerError)
		SemesterCache.Mutex.Unlock()
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

	id, err := strconv.ParseInt(stringid, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	err = db.DeleteSemester(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SemesterCache.Mutex.Lock()
	delete(SemesterCache.Data, id)
	SemesterCache.Mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func PutSemester(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	if stringid == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(stringid, 10, 64)
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

	SemesterCache.Mutex.Lock()
	delete(SemesterCache.Data, id)
	SemesterCache.Mutex.Unlock()

	json.NewEncoder(w).Encode(semester)
}

func GetSemesterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringid := vars["id"]
	if stringid == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(stringid, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SemesterCache.Mutex.RLock()
	if semester, found := SemesterCache.Data[id]; found {
		SemesterCache.Mutex.RUnlock()
		json.NewEncoder(w).Encode(semester)
		return
	}
	SemesterCache.Mutex.RUnlock()

	semester, err := db.GetSemesterByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SemesterCache.Mutex.Lock()
	if semester.ID != nil {
		SemesterCache.Data[*semester.ID] = *semester
	} else {
		http.Error(w, "Semester ID is nil", http.StatusInternalServerError)
	}
	SemesterCache.Mutex.Unlock()

	json.NewEncoder(w).Encode(semester)
}
