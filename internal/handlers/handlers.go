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
	"referral/paradise/internal/tools"

	"github.com/go-sql-driver/mysql"
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

	// Step 1: Validate the referral code with an external request
	resp, err := tools.CheckReferralCode(ref)
	if err != nil {
		log.Printf("Error validating referral code: %v", err)
		http.Error(w, "Failed to validate referral code", http.StatusInternalServerError)
		return
	}

	if !resp {
		http.Error(w, "Invalid referral code", http.StatusBadRequest)
		return
	}

	// Step 2: Attempt to insert the referral code into the database
	_, err = db.DB.Exec("INSERT INTO ref_code (ref) VALUES (?)", ref)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if isDuplicateEntryError(err) {
			http.Error(w, "Referral code already inserted", http.StatusConflict)
		} else {
			log.Printf("Error inserting referral code: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Step 3: Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Referral code successfully added")
}

// isDuplicateEntryError checks if the error is due to a duplicate entry in the database
func isDuplicateEntryError(err error) bool {
	// This logic will depend on your specific SQL driver
	// For example, with MySQL you might check for a specific error number
	// Adjust this to match your SQL driver
	if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1062 {
		return true
	}
	return false
}
