package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	DatabasePath string `env:"DATABASE_PATH" envDefault:"./test.db"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}

	database, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Go SQLite Web App!")
	})

	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, firstName TEXT, lastName TEXT)")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statement.Exec()
		fmt.Fprintf(w, "Table created")
	})

	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		// POST requestではなかったらエラー
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		type User struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		statement, _ := database.Prepare("INSERT INTO users (firstName, lastName) VALUES (?, ?)")
		statement.Exec(user.FirstName, user.LastName)
		fmt.Fprintf(w, "New user was added")
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.Query("SELECT id, firstName, lastName FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type User struct {
			ID        int    `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		var users []User
		for rows.Next() {
			var p User
			rows.Scan(&p.ID, &p.FirstName, &p.LastName)
			users = append(users, p)
		}

		json.NewEncoder(w).Encode(users)
	})

	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
