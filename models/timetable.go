package models

import (
	"encoding/json"
	"time"
)

type Semester struct {
	ID       *int64    `json:"id,omitempty"`
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

	if aux.ID != nil {
		id := int64(*aux.ID)
		s.ID = &id
	} else {
		s.ID = nil
	}
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
