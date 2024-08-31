package main

import (
	"log"
	"net/http"

	"referral/paradise/internal/db"
	"referral/paradise/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()
	defer db.Close()

	db.RunMigrations()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/add", handlers.AddRefHandler).Methods("POST")

	mux := http.NewServeMux()

	mux.Handle("/", r)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8050", mux))
}
