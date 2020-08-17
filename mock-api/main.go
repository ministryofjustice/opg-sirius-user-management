package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Errors string `json:"errors"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth/change-password", ChangePassword).Methods(http.MethodPost, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	existingPassword := r.FormValue("existingPassword")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	errors := ""

	if existingPassword == "" || password == "" || confirmPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		errors = "Missing required field"
	} else if existingPassword != "Password1" {
		w.WriteHeader(http.StatusBadRequest)
		errors = "Password supplied was incorrect or user is not active"
	} else if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		errors = "Confirmation did not match new password"
	}

	response := Response{Errors: errors}

	json.NewEncoder(w).Encode(response)
}
