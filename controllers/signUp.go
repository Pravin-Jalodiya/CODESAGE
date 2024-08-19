package controllers

import (
	"bufio"
	"cli-project/LLM"
	"cli-project/config"
	"cli-project/models"
	passwrd "cli-project/utils/password"
	usr "cli-project/utils/user"
	"cli-project/utils/writers"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

var (
	reader   = bufio.NewReader(os.Stdin)
	username string
	password string
)

func SignUp() {

	for {
		fmt.Print("\nEnter your desired username: ")
		user, err := reader.ReadString('\n')
		user = strings.TrimSuffix(user, "\n")
		user = strings.TrimSpace(user)
		if err != nil {
			fmt.Println("Error reading input.")
			//log errors concurrently
			return
		} else {
			if !usr.IsUnique(user) {
				fmt.Println("Username already exists. Please choose a different username.")
			} else {
				username = user
				break
			}
		}
	}

	for {
		fmt.Print("Enter your password: ")
		secret1, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		pass := string(secret1)
		pass = strings.TrimSpace(pass)
		if err != nil {
			fmt.Println("Error reading input.")
			return
		} else {
			if !passwrd.ValidatePass(pass) {
				fmt.Println("Weak password.\nTry another password (should be at least 8  characters long and must have at least 1 lowercase, 1 uppercase, 1 special and 1 digit characters.)")
				fmt.Println("Need help with finding the right password?(y/n)")
				for {
					choice, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error reading input.")
					}
					choice = strings.TrimSuffix(choice, "\n")
					choice = strings.TrimSpace(choice)
					if choice == "y" {
						suggestedPass := LLM.PasswordSuggestion()
						fmt.Println("Password suggestion : " + suggestedPass)
						break
					} else {
						if choice == "n" {
							break
						}
						fmt.Println("Invalid input. Please try again.")
					}
				}
			} else {
				fmt.Print("Enter your password again: ")
				secret2, err := term.ReadPassword(syscall.Stdin)
				fmt.Println()
				if err != nil {
					fmt.Println("Error reading input.")
					return
				}
				confirmationPass := string(secret2)
				confirmationPass = strings.TrimSpace(confirmationPass)
				if pass == confirmationPass {
					password = confirmationPass
					break
				}
				fmt.Println("Passwords don't match. Try again.")
			}
		}
	}

	hashedPassword, err := passwrd.HashPass(password)
	if err != nil {
		fmt.Println("Error hashing password.")
		return
	}

	email := ""

	newUser := models.User{
		UserID:   1,
		Username: username,
		Password: hashedPassword,
		Name:     username,
		Email:    email,
		Role:     "user",
	}

	ok, err := writers.FWriterUser(config.USER_FILE, newUser)
	if ok {
		fmt.Println("Sign Up successful!")
	} else {
		fmt.Println("Sign up failed : ", err)
	}
}
