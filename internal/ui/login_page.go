package ui

import (
	"cli-project/internal/app/services"
	"cli-project/internal/config/roles"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"cli-project/pkg/validation"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func (ui *UI) ShowLoginPage() {

	// Clear the screen
	fmt.Print("\033[H\033[2J")

	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("               LOGIN                 ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))

	var username, password string
	for {
		// Read Username
		fmt.Print(formatting.Colorize("Username: ", "blue", ""))
		username, _ = ui.reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Printf("%s Username cannot be empty. Try again.\n", emojis.Info)
			continue
		}

		if validation.ValidateUsername(username) {
			fmt.Printf("%s Username is invalid. Try again.\n", emojis.Info)
			continue
		}

		break
	}

	for {
		// Read Password
		fmt.Print(formatting.Colorize("Password: ", "blue", ""))
		passwordBytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password = string(passwordBytes)
		if password == "" {
			fmt.Printf("%s Password cannot be empty. Try again.\n", emojis.Info)
			continue
		}
		break
	}

	// Attempt to log in
	err := ui.userService.Login(username, password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			fmt.Println(emojis.Error, "User not found. Would you like to sign up instead? (y/n)")
			for {
				fmt.Print(formatting.Colorize("Choice: ", "blue", ""))
				choice, err := ui.reader.ReadString('\n')
				choice = strings.TrimSuffix(choice, "\n")
				if err != nil {
					fmt.Println(emojis.Error, "Failed to read input. Please try again.")
					continue
				}
				choice = strings.TrimSpace(choice)
				if strings.ToLower(choice) == "y" {
					ui.ShowSignupPage()
					return
				} else if strings.ToLower(choice) == "n" {
					break
				} else {
					fmt.Println("Invalid input. Please try again.")
				}
			}
		} else if errors.Is(err, services.ErrInvalidCredentials) {
			fmt.Println(emojis.Error, "Username or password incorrect. Please try again.")
		} else {
			fmt.Println(emojis.Error, "Login failed:", err)
		}
	} else {
		fmt.Println(emojis.Success, "Login successful!")

		globals.ActiveUser = username

		role, err := ui.userService.GetUserRole(globals.ActiveUser)
		if err != nil {
			fmt.Println("Unexpected Error:", err)
		}

		if role == roles.USER {
			ui.ShowUserMenu()
		} else if role == roles.ADMIN {
			ui.ShowAdminMenu()
		}

		return
	}
}
