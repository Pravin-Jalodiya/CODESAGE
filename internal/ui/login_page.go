package ui

import (
	"cli-project/internal/app/services"
	"cli-project/internal/config/roles"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
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

	fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
	fmt.Println(utils.Colorize("               LOGIN                 ", "cyan", "bold"))
	fmt.Println(utils.Colorize("====================================", "cyan", "bold"))

	var username, password string
	for {
		// Read Username
		fmt.Print(utils.Colorize("Username: ", "yellow", ""))
		username, _ = ui.reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")
		username = utils.CleanString(username)

		if !validation.ValidateUsername(username) {
			if len(username) == 0 {
				fmt.Printf("%s Username cannot be empty. Try again.\n", utils.InfoEmoji)
			} else {
				fmt.Printf("%s Username is invalid. Try again.\n", utils.InfoEmoji)
			}
			continue
		}

		// Read Password
		fmt.Print(utils.Colorize("Password: ", "yellow", ""))
		passwordBytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password = string(passwordBytes)

		if password == "" {
			fmt.Printf("%s Password cannot be empty. Try again.\n", utils.InfoEmoji)
			continue
		}
		fmt.Println()

		// Attempt to log in
		err := ui.userService.Login(username, password)
		if err != nil {

			var choice string

			if errors.Is(err, services.ErrUserNotFound) {
				fmt.Println(utils.ErrorEmoji, "User not found. Would you like to sign up instead? (y(es)/n(o))")
				for {
					fmt.Print(utils.Colorize("Choice: ", "yellow", ""))
					choice, err = ui.reader.ReadString('\n')
					choice = strings.TrimSuffix(choice, "\n")
					choice = strings.TrimSpace(choice)

					if err != nil {
						fmt.Println(utils.ErrorEmoji, "Failed to read input. Please try again.")
						continue
					}

					if strings.ToLower(choice) == "y" || choice == "yes" {
						ui.ShowSignupPage()
						return

					} else if strings.ToLower(choice) == "n" || choice == "no" {
						break

					} else {
						fmt.Println("Invalid input. Please try again.")
						continue
					}
				}

				continue

			} else if errors.Is(err, services.ErrInvalidCredentials) {
				fmt.Println(utils.ErrorEmoji, "Username or password incorrect. Please try again.")
				continue

			} else {
				fmt.Println(utils.ErrorEmoji, "Login failed:", err)
			}

		} else {
			fmt.Println(utils.SuccessEmoji, "Login successful!")

			globals.ActiveUserID, err = ui.userService.GetUserID(username)

			if err != nil {
				fmt.Println(utils.ErrorEmoji, "Failed to get user ID:", err)
				return
			}

			role, err := ui.userService.GetUserRole(globals.ActiveUserID)
			banned, err := ui.userService.IsUserBanned(globals.ActiveUserID)
			if err != nil {
				fmt.Println("Unexpected Error:", err)
			}

			if banned {
				ui.ShowBannedMessage()
			} else if role == roles.USER {
				ui.ShowUserMenu()
			} else if role == roles.ADMIN {
				ui.ShowAdminMenu()
			}

		}
		return
	}
}
