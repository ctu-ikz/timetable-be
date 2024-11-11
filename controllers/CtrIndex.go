package controllers

import (
	"encoding/json"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"version": "1.0"})
}
