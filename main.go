package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type TimeResponse struct {
	Time string `json:"time"`
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/time_api")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// connection test
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Failed to load Toronto timezone", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().In(loc)
	timeString := currentTime.Format("2006-01-02 15:04:05")

	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", timeString)
	if err != nil {
		http.Error(w, "Failed to log time to database", http.StatusInternalServerError)
		return
	}
	response := TimeResponse{Time: timeString}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/current-time", currentTimeHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
