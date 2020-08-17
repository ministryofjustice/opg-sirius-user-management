package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/change-password", changePassword)

	log.Println("Running at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}

type errorsResponse struct {
	Errors string `json:"errors"`
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	existingPassword := r.FormValue("existingPassword")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	errorMessage, ok := validate(existingPassword, password, confirmPassword)

	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(errorsResponse{Errors: errorMessage})
}

func validate(existingPassword, password, confirmPassword string) (string, bool) {
	if existingPassword == "" || password == "" || confirmPassword == "" {
		return "Missing required field", false
	} else if existingPassword != "Password1" {
		return "Password supplied was incorrect or user is not active", false
	} else if password != confirmPassword {
		return "Confirmation did not match new password", false
	}

	return "", true
}
