package ui

import (
	"cli-project/pkg/globals"
	"cli-project/pkg/utils/formatting"
	"fmt"
)

func (ui *UI) ShowUserProfile() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Fetch the user profile details (assuming `ui.userService.GetUserProfile` returns the user's profile)
	user, err := ui.userService.GetUserByID(globals.ActiveUserID)
	if err != nil {
		fmt.Println(formatting.Colorize("Failed to load user profile.", "red", "bold"))
		return
	}

	// Display the user profile
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("            USER PROFILE            ", "cyan", "bold"))
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("Username: ", "magenta", "bold"), user.StandardUser.Username)
	fmt.Println(formatting.Colorize("Name: ", "magenta", "bold"), user.StandardUser.Name)
	fmt.Println(formatting.Colorize("Email: ", "magenta", "bold"), user.StandardUser.Email)
	fmt.Println(formatting.Colorize("Leetcode ID: ", "magenta", "bold"), user.LeetcodeID)
	fmt.Println(formatting.Colorize("Organisation: ", "magenta", "bold"), user.StandardUser.Organisation)
	fmt.Println(formatting.Colorize("Country: ", "magenta", "bold"), user.StandardUser.Country)

	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("Press any key to return to the main menu.", "", "bold"))

	// Wait for the user to press any key to return
	ui.reader.ReadString('\n')
}
