package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	_ "github.com/mattn/go-sqlite3"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var mu sync.Mutex

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if !app.isValidCredentials(username, password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users := app.getUsers()

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) isValidCredentials(username, password string) bool {
	return username == "JIIL" && password == "juice"
}

func (app *application) getUsers() []auth {
	mu.Lock()
	defer mu.Unlock()

	rows, err := app.db.Query("SELECT username, password FROM users")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var users []auth
	for rows.Next() {
		var user auth
		if err := rows.Scan(&user.Username, &user.Password); err != nil {
			log.Println(err)
			return nil
		}
		users = append(users, user)
	}

	return users
}

func (app *application) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user auth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, err := app.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User added successfully")
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user auth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	result, err := app.db.Exec("DELETE FROM users WHERE username = ? AND password = ?", user.Username, cloud.HashPassword(user.Password))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	fmt.Fprint(w, "User deleted successfully")
}
