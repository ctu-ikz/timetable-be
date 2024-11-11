package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(fmt.Sprintf("Welcome to the timetable-be API \nVersion: 1.0"))
}
