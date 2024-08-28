package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
    "math/rand"


	"referral/paradise/internal/db"
)


func generateSeed(s string) int64 {
	seed := int64(0)
	for _, char := range s {
		seed += int64(char)
	}
	return seed
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get IP Address
	ipAddress := r.RemoteAddr

	// Generate seed
	salt := os.Getenv("SALT")
	now := time.Now()
	seedStr := fmt.Sprintf("%s%s%d%d", salt, ipAddress, now.Day(), now.Year())
	seed := generateSeed(seedStr)

	rand.Seed(seed)

	// Get a random number from seed
	randomNumber := rand.Int()

	// Get the table size
	tableSize, err := db.GetTableSize()
	if err != nil {
		log.Printf("Error fetching table size: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if tableSize < 1 {
		tableSize = 1
	}

	// Calculate the row number
	rowNumber := randomNumber % tableSize

	// Fetch the ref string from the table
	var ref string
	err = db.DB.QueryRow("SELECT ref FROM ref_code LIMIT 1 OFFSET ?", rowNumber).Scan(&ref)
	if err != nil {
		log.Printf("Error fetching ref: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the HTML template
	tmpl, err := template.New("home").Parse("<html><body><h1>{{.}}</h1></body></html>")
	if err != nil {
		log.Printf("Error creating template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, ref)
}

func AddRefHandler(w http.ResponseWriter, r *http.Request) {
	ref := r.FormValue("ref")
	if ref == "" {
		http.Error(w, "Missing ref parameter", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("INSERT INTO ref_code (ref) VALUES (?)", ref)
	if err != nil {
		log.Printf("Error inserting ref: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Ref added successfully")
}
