package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "123"
	dbName   = "gocrud_app"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/user/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", deleteUserHandler).Methods("DELETE")

	log.Println("server listening on 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	CreateUser(db, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed Create User", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User Create Success")
}

func CreateUser(db *sql.DB, name, email string) error {
	query := "INSERT INTO users (name,email) VALUES (?,?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	ID    int
	Name  string
	Email string
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)

	user, err := GetUser(db, userID)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func GetUser(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)

	UpdateUser(db, userID, user.Name, user.Email)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusFound)
		return
	}
	fmt.Fprintln(w, "User Update Success")
}
func UpdateUser(db *sql.DB, id int, name, email string) (*User, error) {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	_, err := db.Exec(query, name, email, id)

	if err != nil {
		return nil, err
	}
	return &User{}, nil
}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	user := DeleteUser(db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")

	// Convert the user object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func DeleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
