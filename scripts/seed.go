package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Helper to generate a random minute offset
func randomMinutes(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min+1)+min) * time.Minute
}

func main() {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Failed to load Asia/Kolkata: %v. Falling back to UTC.", err)
	} else {
		time.Local = loc
	}

	db, err := sql.Open("sqlite3", "./events.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Clear existing data to avoid conflicts for testing
	_, err = db.Exec("DELETE FROM events")
	if err != nil {
		log.Fatal("Failed to clear existing events:", err)
	}
	fmt.Println("Cleared old events.")

	// Start from 30 days ago
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)

	// We seed logic day by day
	// Roughly normal work day:
	// ~08:00 AM leave home -> outside
	// ~08:30 AM arrive at office -> office
	// ~12:30 PM lunch break -> outside
	// ~01:15 PM return -> office
	// ~05:00 PM leave office -> outside
	// ~05:45 PM arrive home -> home
	// Sometimes maybe they wfh (home all day, skip logic)
	// Sometimes weekend (home -> outside -> home)

	events := 0
	for d := 0; d <= 30; d++ {
		currentDate := startDate.AddDate(0, 0, d)
		weekday := currentDate.Weekday()

		insert := func(place string, t time.Time) {
			_, err := db.Exec("INSERT INTO events(place, timestamp) VALUES(?, ?)", place, t)
			if err != nil {
				log.Fatalf("Insert failed: %v", err)
			}
			events++
		}

		// Calculate start of day
		baseTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
			0, 0, 0, 0, currentDate.Location())

		// Weekend pattern
		if weekday == time.Saturday || weekday == time.Sunday {
			// Leave home around 11am
			leaveTime := baseTime.Add(11 * time.Hour).Add(randomMinutes(-60, 60))
			insert("outside", leaveTime)

			// Return home around 4pm
			returnTime := leaveTime.Add(5 * time.Hour).Add(randomMinutes(-60, 60))
			insert("home", returnTime)
			continue
		}

		// Workday pattern (with small 10% chance of WFH where we just stay at home)
		if rand.Float32() < 0.10 {
			continue // WFH, no logged transitions (assumes implicitly at home)
		}

		// Leave home for commute (around 8am)
		leaveHomeTime := baseTime.Add(8 * time.Hour).Add(randomMinutes(-30, 30))
		insert("outside", leaveHomeTime)

		// Arrive at office
		arriveOfficeTime := leaveHomeTime.Add(randomMinutes(20, 45))
		insert("office", arriveOfficeTime)

		// Lunch out (around 12:30pm)
		lunchOutTime := baseTime.Add(12*time.Hour + 30*time.Minute).Add(randomMinutes(-15, 30))
		insert("outside", lunchOutTime)

		// Return from lunch
		lunchReturnTime := lunchOutTime.Add(randomMinutes(30, 60))
		insert("office", lunchReturnTime)

		// Leave office (around 5pm)
		leaveOfficeTime := baseTime.Add(17 * time.Hour).Add(randomMinutes(-30, 60))
		insert("outside", leaveOfficeTime)

		// Arrive home
		arriveHomeTime := leaveOfficeTime.Add(randomMinutes(25, 50))
		insert("home", arriveHomeTime)
	}

	fmt.Printf("Seed complete! Inserted %d realistic logical events over the past 30 days.\n", events)
}
