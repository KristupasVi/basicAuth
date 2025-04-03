package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	username string
	email    string
	password string
}

func main() {
	InitDB()
	start()
}

func start() {
	defer DB.Close()
	for {
		fmt.Println("\nChoose action:")
		fmt.Println("1. Register")
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
	u := user{}

	// Username logic
	u.username = fieldScanVal("Enter your username: ", usernameValidation)

	var existsUsername bool
	var err error

	if _, existsUsername, err = getUserByName(u.username); err != nil {
		fmt.Printf("Error checking username %v\n", err)
		return
	} else if existsUsername {
		fmt.Printf("This username \"%s\" already taken.\n", u.username)
		return
	}

	// Email logic
	u.email = fieldScanVal("Enter your email: ", emailValidation)

	var existsEmail bool

	if _, existsEmail, err = getUserByEmail(u.email); err != nil {
		fmt.Printf("Error checking email %v\n", err)
	} else if existsEmail {
		fmt.Printf("This email \"%s\" already exists", u.email)
		return
	}

	// Password logic
	if !existsUsername && !existsEmail {
		password := fieldScanVal("Enter your password: ", passwordValidation)
		// Pasword hashing
		hashedPassword, err := passwordHash(password)
		if err != nil {
			fmt.Printf("An error occurred while hashing the password: %v\n", err)
			return
		}
		u.password = hashedPassword
		// If everything is successful
		err = registerUser(u.username, u.email, u.password)
		if err != nil {
			fmt.Printf("Error registering user: %v\n", err)
			return
		}
		fmt.Println("User created successfully")
	}
}

func login() {

	var attempts int
	const maxAttempts = 3

	for i := 0; i < maxAttempts; i++ {

		// Input fields
		identifier := fieldScanVal("Enter your username or email: ", nil)
		password := fieldScanVal("Enter your password: ", nil)

		var hashedPassword string
		var exists bool
		var err error

		// Determines whether the input is an email or username
		if strings.Contains(identifier, "@") {
			hashedPassword, exists, err = getUserByEmail(identifier)
		} else {
			hashedPassword, exists, err = getUserByName(identifier)
		}
		if err != nil {
			fmt.Printf("Error checking credentials: %v\n", err)
			return
		}

		if !exists {
			fmt.Println("Invalid username/email or password")
		} else if !passwordVerification(hashedPassword, password) {
			fmt.Println("Invalid username/email or password")
		} else {
			fmt.Println("Login successful!")
			return
		}
		attempts++
	}
	fmt.Println("Too many failed attempts. Please try again later")
}

// Hashing user passwordInput register stage
func passwordHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verifying user password input with hashedPassword login stage
func passwordVerification(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Checks if the data provided was empty string in a field
func checkIfBlank(field, message string) string {
	for field == "" {
		fmt.Println("This field cannot be empty")
		fmt.Print(message)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		field = strings.TrimSpace(scanner.Text())
	}
	return field
}

// Scans a field Input and validates it
func fieldScanVal(fieldText string, validation func(string) (string, bool)) string {
	fmt.Print(fieldText)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fieldInput := strings.TrimSpace(scanner.Text())
	fieldInput = checkIfBlank(fieldInput, fieldText)
	if validation != nil {
		if message, correct := validation(fieldInput); !correct {
			fmt.Println(message)
			return fieldScanVal(fieldText, validation)
		}
	}
	return fieldInput
}
