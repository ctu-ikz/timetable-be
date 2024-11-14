package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ctu-ikz/timetable-be/helpers"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer", "", 1))
		err := helpers.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Errorf("Invalid token").Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
