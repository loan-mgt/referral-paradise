package handlers

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"referral/paradise/internal/db"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

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

	// Remove the port from the IP address
	ipAddress = strings.Split(ipAddress, ":")[0]

	// Generate seed
	salt := os.Getenv("SALT")
	now := time.Now()
	seedStr := fmt.Sprintf("%s%s%d%d", salt, ipAddress, now.Day(), now.Year())
	seed := generateSeed(seedStr)

	random := rand.New(rand.NewSource(seed))

	// Get a random number from seed
	randomNumber := random.Int()

	log.Printf("rand seed %d number %d seedStr %s", seed, randomNumber, seedStr)
	//mysalt127.0.0.1:62931282024

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
	err = templates.ExecuteTemplate(w, "home.html", ref)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func AddRefHandler(w http.ResponseWriter, r *http.Request) {
	ref := r.FormValue("ref")
	if ref == "" {
		http.Error(w, "Missing ref parameter", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("INSERT INTO ref_code (ref) VALUES (?)", ref)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
