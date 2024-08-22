package ui

import (
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"cli-project/pkg/validation"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func (ui *UI) ShowLoginPage() {
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("               LOGIN                 ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))

	for {
		// Read Username
		var username string
		fmt.Print(formatting.Colorize("Username: ", "blue", ""))
		username, _ = ui.reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Printf("%s Username cannot be empty. Try again.\n", emojis.Info)
			continue
		}
		if validation.ValidateUsername(username) {

		}

		// Read Password
		fmt.Print(formatting.Colorize("Password: ", "blue", ""))
		passwordBytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password := string(passwordBytes)
		if password == "" {
			fmt.Printf("%s Password cannot be empty. Try again.\n", emojis.Info)
			continue
		}
		fmt.Println()

		// Attempt to login
		err := ui.userService.Login(username, password)
		if err != nil {
			if err.Error() == "user not found" {
				fmt.Println(emojis.Error, "User not found. Would you like to sign up instead? (y/n)")
				var choice string
				fmt.Print(formatting.Colorize("Choice: ", "blue", ""))
				choice, err := ui.reader.ReadString('\n')
				if err != nil {
					fmt.Println(emojis.Error, "Failed to read input. Please try again.")
					return
				}
				choice = strings.TrimSpace(choice)
				if strings.ToLower(choice) == "y" {
					ui.ShowSignupPage()
					return
				}
			} else if err.Error() == "username or password incorrect" {
				fmt.Println(emojis.Error, "Username or password incorrect. Please try again.")
			} else {
				fmt.Println(emojis.Error, "Login failed:", err)
			}
		} else {
			fmt.Println(emojis.Success, "Login successful!")
			return
		}
	}
}
