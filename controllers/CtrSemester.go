package controllers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
)

var (
	cache      *models.SemesterCache
	cacheMutex sync.RWMutex
)

func GetDbSemester(w http.ResponseWriter, r *http.Request) {
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
