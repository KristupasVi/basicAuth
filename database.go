// database.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Initially it's nil
var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	psqlConnnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnnection)
	if err != nil {
		log.Fatal(err)
	}
	// We di ping to test if the database is reachable
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	createUserTable()
	fmt.Println("Successfully connected to PostgreSQL!")
}

func createUserTable() {
	query := `
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			username VARCHAR(20) UNIQUE NOT NULL,
			email VARCHAR(70) UNIQUE NOT NULL,
			password VARCHAR(70) NOT NULL

		)
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func registerUser(username, email, hashedPassword string) error {
	query := "INSERT INTO users (username,email, password) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, username, email, hashedPassword)
	return err
}

func getUserByName(username string) (string, bool, error) {
	var hashedPassword string
	// QueryRow doesn't give me any data it prepares the row and waits for extraction .Scan is how the data is pulled out of that row
	err := DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	// It can be like this : row := DB.QueryRow() err := row.Scan() if there is a need to do something specific with a row
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return hashedPassword, true, nil
}

func getUserByEmail(email string) (string, bool, error) {
	var hashedPassword string

	err := DB.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return hashedPassword, true, nil
}
