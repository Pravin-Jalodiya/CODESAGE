package controllers

import (
	"cli-project/UI"
	"cli-project/config/roles"
	"cli-project/middleware"
	passwrd "cli-project/utils/password"
	"cli-project/utils/readers"
	"cli-project/utils/str"
	usr "cli-project/utils/user"
	"fmt"
	"golang.org/x/term"
	"strings"
	"syscall"
)

func Login() {

	fmt.Print("Enter your username: ")
	user, err2 := reader.ReadString('\n')
	user = strings.TrimSuffix(user, "\n")
	user = strings.TrimSpace(user)
	if err2 != nil {
		fmt.Println("Error reading input.")
		//log errors
		return
	} else {
		username = user
	}

	fmt.Print("Enter your password: ")
	secret, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println("Error reading input.")
		return
	}
	fmt.Println()
	pass := string(secret)
	pass = strings.TrimSpace(pass)
	password = pass

	if storedHash, exists := readers.UserPassMap[username]; exists {
		if passwrd.VerifyPassword(password, storedHash) {
			fmt.Println(str.Colorize("Log In successful!", "green", "", "bold"))
			middleware.Auth(usr.GetUser(username).UserID)
			if middleware.VerifyRole(username, roles.USER) {
				UI.User()
			} else if middleware.VerifyRole(username, roles.ADMIN) {
				//UI.admin()
			}
		} else {
			fmt.Println("Incorrect username or password")
		}
	} else {
		fmt.Println("User not found. Do you wish to sign up?(y/n)")
		signupChoice, err := reader.ReadString('\n')
		signupChoice = strings.TrimSuffix(signupChoice, "\n")
		signupChoice = strings.TrimSpace(signupChoice)
		if err != nil {
			fmt.Println("Error reading input.")
			return
		}
		if signupChoice == "y" {
			SignUp()
		} else {
			return
		}

	}

}
