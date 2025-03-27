package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string) // key: username, value: password

func main() {
	start()
}

func start() {
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

	// This is map lookup  we ignore the value with _ if the is a value then exists = true and we execute the code {}
	if _, exists := users[username]; exists {
		fmt.Printf("The username \"%s\" is already taken.\n", username)
	} else {
		fmt.Print("Write your password: ")
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())
		// Pasword hashing
		hashedPassword, err := passwordHash(password)
		if err != nil {
			fmt.Printf("An error occurred while hashing the password: %v\n", err)
			return
		}
		users[username] = hashedPassword
		//fmt.Println("HashedPassword", hashedPassword)
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

		if hashedPassword, exists := users[username]; exists {
			if passwordVerification(hashedPassword, password) {
				fmt.Println("Login was successful")
				break
			} else {
				attempts++
				fmt.Printf("Wrong password! You have %d/3 attempts left\n", attempts)
			}
		} else {
			fmt.Println("User doesn't exist")
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
