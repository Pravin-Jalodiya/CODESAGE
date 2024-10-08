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
	fmt.Println(formatting.Colorize("Username: ", "cyan", "bold"), user.StandardUser.Username)
	fmt.Println(formatting.Colorize("Name: ", "cyan", "bold"), user.StandardUser.Name)
	fmt.Println(formatting.Colorize("Email: ", "cyan", "bold"), user.StandardUser.Email)
	fmt.Println(formatting.Colorize("Leetcode ID: ", "cyan", "bold"), user.LeetcodeID)
	fmt.Println(formatting.Colorize("Organisation: ", "cyan", "bold"), user.StandardUser.Organisation)
	fmt.Println(formatting.Colorize("Country: ", "cyan", "bold"), user.StandardUser.Country)

	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')

}
