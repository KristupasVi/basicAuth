package main

import "regexp"

func usernameValidation(username string) (returnMessage string, correctInput bool) {

	if len(username) < 6 || len(username) > 25 {
		return "Username should be between 6 and 25 characters long", false
	}
	re := regexp.MustCompile(`^.[a-zA-Z0-9]+$`)
	if !re.MatchString(username) {
		return "Username can only contain letters and numbers", false
	}
	return "", true
}

func emailValidation(email string) (returnMessage string, correctInput bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$")
	if !re.MatchString(email) {
		return "This email format is invalid ", false
	}
	return "", true
}

func passwordValidation(password string) (returnMessage string, correctInput bool) {
	if len(password) < 6 || len(password) > 30 {
		return "Password should contain atleast 1 number, 1 symbol and the length should be atleast 6 characters ", false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return "Password should contain atleast 1 number, 1 symbol and the length should be atleast 6 characters ", false
	}
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
	if !hasSpecial {
		return "Password should contain atleast 1 number, 1 symbol and the length should be atleast 6 characters ", false
	}
	return "", true
}
