package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/helpers"
	"github.com/ctu-ikz/timetable-be/models"
	"github.com/gorilla/mux"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == nil || *user.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	var dbUser *models.User
	dbUser, err = db.GetUserByUsername(user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dbUser == nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !helpers.CheckPassword(*user.Password, *dbUser.Password) {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	jwt := "tady bude token"

	json.NewEncoder(w).Encode(jwt)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userNameTaken, err := db.UserNameTaken(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userNameTaken {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	newUser, err := db.PostUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newUser)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
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

	user, err := db.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
