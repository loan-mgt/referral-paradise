package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "referral/paradise/internal/db"
    "referral/paradise/internal/handlers"
)

func main() {
    db.Connect()
    defer db.Close()

    db.RunMigrations()

    r := mux.NewRouter()
    r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
    r.HandleFunc("/add", handlers.AddRefHandler).Methods("POST")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8050", nil))
}
