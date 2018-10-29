package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var YourHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
})

// Get single book
var getBook = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r) // Gets params
	json.NewEncoder(w).Encode(mux.Vars(r))
})
