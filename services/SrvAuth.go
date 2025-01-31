package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/helpers"
	"github.com/ctu-ikz/timetable-be/models"
	"github.com/gorilla/mux"
)

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserDB
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := db.GetUserByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dbUser == nil || dbUser.ID == nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !helpers.CheckPassword(*user.Password, *dbUser.Password) {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := helpers.CreateTokens(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokenSeconds, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		http.Error(w, "Invalid refresh token expiry duration", http.StatusInternalServerError)
		return
	}

	hashedRefreshToken := HashToken(refreshToken)
	if err := db.CreateRefreshToken(*dbUser.ID, hashedRefreshToken, time.Now().Add(time.Duration(refreshTokenSeconds)*time.Second)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserDB
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

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
		return
	}

	refreshTokenString := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer ", "", 1))

	if err := helpers.VerifyToken(refreshTokenString, true); err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	userID, err := helpers.GetUserIDFromToken(refreshTokenString, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "User ID parsing error", http.StatusUnauthorized)
		return
	}

	hashedRefreshTokenOld := HashToken(refreshTokenString)

	validToken, err := db.CheckRefreshTokenValidity(hashedRefreshTokenOld, userID64)
	if err != nil || !validToken {
		http.Error(w, "Invalid refresh token in database", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(userID64)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := db.InvalidateRefreshToken(hashedRefreshTokenOld); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newAccessToken, newRefreshToken, err := helpers.CreateTokens(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedRefreshToken := HashToken(newRefreshToken)

	if err := db.CreateRefreshToken(userID64, hashedRefreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}
