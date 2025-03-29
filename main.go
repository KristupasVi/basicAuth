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

	// Username Logic
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	u.username = strings.TrimSpace(scanner.Text())

	checkIfBlank(u.username, "Enter your username: ")

	var existsUsername bool
	var err error

	if _, existsUsername, err = getUserByName(u.username); err != nil {
		fmt.Printf("Error checking username %v\n", err)
		return
	} else if existsUsername {
		fmt.Printf("The username \"%s\" is already taken.\n", u.username)
		return
	}

	// Email Logic
	fmt.Print("Enter your email: ")
	scanner.Scan()
	u.email = strings.TrimSpace(scanner.Text())

	checkIfBlank(u.email, "Enter your email: ")

	var existsEmail bool

	if _, existsEmail, err = getUserByEmail(u.email); err != nil {
		fmt.Printf("Error checking email %v\n", err)
	} else if existsEmail {
		fmt.Printf("The email \"%s\" is already exists", u.email)
	}

	// Password logic
	if !existsUsername && !existsEmail {
		fmt.Print("Enter your password: ")
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())

		checkIfBlank(password, "Enter your password: ")

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

	fmt.Print("Enter your username or email: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	identifier := strings.TrimSpace(scanner.Text())

	checkIfBlank(identifier, "Enter your username or email: ")

	for i := attempts; i < maxAttempts; i++ {
		fmt.Print("Enter your password: ")
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())

		if password == "" {
			fmt.Println("Password cannot be empty")
			continue
		}

		//

		var hashedPassword string
		var exists bool
		var err error

		if strings.Contains(identifier, "@") {
			hashedPassword, exists, err = getUserByEmail(identifier)
		} else {
			hashedPassword, exists, err = getUserByName(identifier)
		}

		if err != nil {
			fmt.Printf("Error checking credentials: %v\n", err)
			return
		}

		if !exists && strings.Contains(identifier, "@") {
			fmt.Println("Invalid email or password try login again")
			return
		} else if !exists {
			fmt.Println("Invalid username or password try login again or use email")
			return
		}

		if passwordVerification(hashedPassword, password) {
			fmt.Println("Login successful!")
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
