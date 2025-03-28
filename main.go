package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	InitDB()
	start()
}

func start() {
	defer DB.Close()
	for {
		fmt.Println("\nChoose action:")
		fmt.Println("1. Registrer")
		fmt.Println("2. Login")
		fmt.Println("3. Leave")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			register()
		case "2":
			login()
		case "3":
			fmt.Println("Bye Bye")
			return
		default:
			fmt.Println("Wrong input. Try again. Available inputs: 1, 2, or 3.")
		}
	}
}

func register() {
	fmt.Print("Write your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if _, exists, err := getUser(username); err != nil {
		fmt.Printf("Error checking username %v\n", err)
		return
	} else if exists {
		fmt.Printf("The username \"%s\" is already taken.\n", username)
		return
	} else { // If username doesn't exists
		fmt.Print("Write your password: ")
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())
		// Pasword hashing
		hashedPassword, err := passwordHash(password)
		if err != nil {
			fmt.Printf("An error occurred while hashing the password: %v\n", err)
			return
		}

		err = registerUser(username, hashedPassword)
		if err != nil {
			fmt.Printf("Error registering user: %v\n", err)
			return
		}
		fmt.Println("User created successfully")
	}
}

func login() {
	var attempts int
	fmt.Print("Write your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	for i := 0; i < 3; i++ {
		fmt.Print("Write your password: ")
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())

		hashedPassword, exists, err := getUser(username)
		if err != nil {
			fmt.Printf("Error checking user: %v\n", err)
			return
		}
		if !exists {
			fmt.Println("User doesn't exist")
			return
		}
		if passwordVerification(hashedPassword, password) {
			fmt.Println("Login was successful")
			return
		} else {
			attempts++
			fmt.Printf("Wrong password! You have %d/3 attempts left\n", attempts)
		}

	}
	if attempts == 3 {
		fmt.Println("Too many failed attempts. Please try again later")
	}
}

// Hashing user passwordInput register stage
func passwordHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	//fmt.Println("Hashed in passwordHash", hashed)
	return string(hashed), nil
}

// Verifying user password input with hashedPassword login stage
func passwordVerification(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
