package models

import (
	"sync"
)

type SemesterCache struct {
	Data  map[int64]Semester
	Mutex sync.RWMutex
}

type TimetableCache struct {
	Data  map[string]WeeklyTimetable
	Mutex sync.RWMutex
}
