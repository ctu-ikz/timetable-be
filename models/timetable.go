package models

import (
	"encoding/json"
	"time"
)

type Semester struct {
	ID       *int      `json:"id,omitempty"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Codename string    `json:"codename"`
}

func (s *Semester) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID       *int   `json:"id,omitempty"`
		Start    string `json:"start"`
		End      string `json:"end"`
		Codename string `json:"codename"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	layout := "2006-01-02"
	var err error
	s.Start, err = time.Parse(layout, aux.Start)
	if err != nil {
		return err
	}

	s.End, err = time.Parse(layout, aux.End)
	if err != nil {
		return err
	}

	s.ID = aux.ID
	s.Codename = aux.Codename

	return nil
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
