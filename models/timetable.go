package models

import "time"

type Semester struct {
	ID       int       `json:"id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Codename string    `json:"codename"`
}

type SubjectClass struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Shortcut  string    `json:"shortcut"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	CodeName  string    `json:"code_name"`
	Day       *int      `json:"day,omitempty"`
	Weeks     *[]string `json:"weeks,omitempty"`
}

type WeeklyTimetable struct {
	Monday    []SubjectClass `json:"monday"`
	Tuesday   []SubjectClass `json:"tuesday"`
	Wednesday []SubjectClass `json:"wednesday"`
	Thursday  []SubjectClass `json:"thursday"`
	Friday    []SubjectClass `json:"friday"`
}
