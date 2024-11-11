package controllers

import (
	"fmt"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Version 1.0")
}
