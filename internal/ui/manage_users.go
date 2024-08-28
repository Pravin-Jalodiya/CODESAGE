package ui

import (
	"cli-project/pkg/utils"
	"cli-project/pkg/utils/data_cleaning"
	"cli-project/pkg/utils/formatting"
	"cli-project/pkg/validation"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

func (ui *UI) ManageUsers() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("           MANAGE USERS             ", "cyan", "bold"))
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("1. View all users", "", ""))
		fmt.Println(formatting.Colorize("2. Ban a user", "", ""))
		fmt.Println(formatting.Colorize("3. Unban a user", "", ""))
		fmt.Println(formatting.Colorize("4. Go back", "", ""))

		fmt.Print(formatting.Colorize("Enter your choice: ", "yellow", "bold"))
		choice, err := ui.reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if err != nil {
			fmt.Println(formatting.Colorize("error reading input:", "red", "bold"), err)
			return
		}

		switch choice {
		case "1":
			ui.viewAllUsers()
		case "2":
			ui.banUser()
		case "3":
			ui.unbanUser()
		case "4":
			return // Go back to the previous menu
		default:
			fmt.Println(formatting.Colorize("invalid choice", "red", "bold"))
		}
	}
}

func (ui *UI) viewAllUsers() {
	// Get all users from the user service
	users, err := ui.userService.GetAllUsers()
	if err != nil {
		fmt.Println("Failed to load users.")
		return
	}

	// If no users found, notify the admin
	if len(*users) == 0 {
		fmt.Println("No users found.")
		return
	}

	// Create a new table writer to format the output as a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Username", "Name", "Email", "Leetcode ID", "Organisation", "Country", "IsBlocked", "Last Seen (IST)"})

	// Print table rows, excluding admin users
	for _, user := range *users {
		if user.StandardUser.Role != "user" {
			continue
		}

		// Convert Last Seen time to IST and format it
		lastSeenIST := utils.ConvertToIST(user.LastSeen)

		// Add the row to the table
		table.Append([]string{
			user.StandardUser.Username,
			user.StandardUser.Name,
			user.StandardUser.Email,
			user.LeetcodeID,
			user.StandardUser.Organisation,
			user.StandardUser.Country,
			fmt.Sprintf("%t", user.StandardUser.IsBanned),
			lastSeenIST,
		})
	}

	// Render the table to the console
	table.Render()

}

func (ui *UI) banUser() {

	// View all users
	ui.viewAllUsers()

	// Logic to ban a user
	var username string
	var err error
	for {
		fmt.Print("Enter the username to ban: ")
		username, err = ui.reader.ReadString('\n')
		username = data_cleaning.CleanString(username)
		if err != nil {
			fmt.Println(formatting.Colorize("error reading input:", "red", "bold"), err)
			return
		}

		valid := validation.ValidateUsername(username)

		if !valid {
			fmt.Println(formatting.Colorize("enter a valid username", "yellow", "bold"))
			continue
		}
		break
	}

	// banning logic
	alreadyBanned, err := ui.userService.BanUser(username)
	if err != nil {
		fmt.Println(formatting.Colorize("user does not exist", "red", "bold"))
		return
	} else if alreadyBanned {
		fmt.Println(formatting.Colorize("user already banned", "yellow", "bold"))
	} else {
		fmt.Println(formatting.Colorize("user banned successfully", "green", "bold"))
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')

}

func (ui *UI) unbanUser() {

	// View all users
	ui.viewAllUsers()

	// Logic to unban a user
	var username string
	var err error
	for {
		fmt.Print("Enter the username to unban: ")
		username, err = ui.reader.ReadString('\n')
		username = data_cleaning.CleanString(username)
		if err != nil {
			fmt.Println(formatting.Colorize("error reading input:", "red", "bold"), err)
			return
		}

		valid := validation.ValidateUsername(username)

		if !valid {
			fmt.Println(formatting.Colorize("enter a valid username", "yellow", "bold"))
			continue
		}
		break
	}

	// Unbanning logic
	alreadyUnbanned, err := ui.userService.UnbanUser(username)
	if err != nil {
		fmt.Println(formatting.Colorize("user does not exist", "red", "bold"), err)
		return
	} else if alreadyUnbanned {
		fmt.Println(formatting.Colorize("user already unbanned", "yellow", "bold"))
	} else {
		fmt.Println(formatting.Colorize("user unbanned successfully", "green", "bold"))
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')

}
