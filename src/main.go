package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func main() {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Health check passed")
	})
	http.HandleFunc("/dbconnection", dbConnectionHandler)

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func dbConnectionHandler(w http.ResponseWriter, r *http.Request) {
	rdsURL := os.Getenv("RDS_URL")
	if rdsURL == "" {
		fmt.Fprint(w, "RDS_URL environment variable is not set")
		return
	}

	// Assuming RDS_URL is in the format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", rdsURL)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Fprintf(w, "Error connecting to database: %s", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Fprintf(w, "Database connection status: DOWN")
	} else {
		fmt.Fprint(w, "Database connection status: UP")
	}
}