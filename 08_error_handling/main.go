package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func ValidateEmail(email string) error {
	rndNumber := rand.Intn(100)
	if rndNumber < 20 {
		return fmt.Errorf("Email already exists: %s", email)
	}
	return nil
}

func ValidatePassword(password string) error {
	rndNumber := rand.Intn(100)
	if rndNumber < 80 {
		return fmt.Errorf("Password does not meet criteria")
	}
	return nil
}

func SaveUser() error {
	rndNumber := rand.Intn(100)
	if rndNumber < 10 {
		return fmt.Errorf("User couldn't be created, try again later")
	}
	return nil
}

func Signup(email, password string) error {
	var validationErrs error
	emailErr := ValidateEmail(email)
	if emailErr != nil {
		validationErrs = errors.Join(validationErrs, fmt.Errorf("Email validation failed: %w", emailErr))
	}
	passErr := ValidatePassword(password)
	if passErr != nil {
		validationErrs = errors.Join(validationErrs, fmt.Errorf("Password validation failed: %w", passErr))
	}
	if validationErrs != nil {
		return fmt.Errorf("Signup validations failed: %w", validationErrs)
	}
	userErr := SaveUser()
	if userErr != nil {
		return fmt.Errorf("SaveUser operation failed: %w", userErr)
	}
	return nil
}

func main() {
	email := "email@domain.com"
	password := "helloworld"
	err := Signup(email, password)
	if err != nil {
		fmt.Printf("Signup operation failed: %v", err)
	} else {
		fmt.Printf("User created successfully")
	}
}
