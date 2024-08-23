package ui

import (
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strings"
)

func (ui *UI) ShowAdminMenu() {
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("             ADMIN MENU             ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("1. View Dashboard", "green", ""))
	fmt.Println(formatting.Colorize("2. Add or Remove Questions", "green", ""))
	fmt.Println(formatting.Colorize("3. Ban or Unban Users", "green", ""))
	fmt.Println(formatting.Colorize("4. Post Announcement", "green", ""))

	fmt.Print(formatting.Colorize("Enter your choice: ", "yellow", "bold"))
	choice, err := ui.reader.ReadString('\n')
	if err != nil {
		fmt.Println(formatting.Colorize("Error reading input:", "red", "bold"), err)
		return
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		ui.ShowAdminDashboard()
	case "2":
		ui.ManageQuestions()
	case "3":
		ui.ManageUsers()
	//case "4":
	//	ui.PostAnnouncement()
	default:
		fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		ui.ShowAdminMenu()
	}
}
