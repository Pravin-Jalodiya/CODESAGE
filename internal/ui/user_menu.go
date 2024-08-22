package ui

import (
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strings"
)

func (ui *UI) ShowUserMenu() {
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("              USER MENU             ", "cyan", "bold"))
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("1. Solve Questions", "green", ""))
	fmt.Println(formatting.Colorize("2. View Dashboard", "green", ""))

	fmt.Print(formatting.Colorize("Enter your choice: ", "yellow", "bold"))
	choice, err := ui.reader.ReadString('\n')
	if err != nil {
		fmt.Println(formatting.Colorize("Error reading input:", "red", "bold"), err)
		return
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		//ui.SolveQuestions()
	case "2":
		fmt.Println(formatting.Colorize("You chose to view the dashboard.", "blue", "bold"))
		//ui.ViewDashboard()
	default:
		fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		ui.ShowUserMenu()
	}
}

//func (u *UI) SolveQuestions() {
//	// Placeholder for solving questions logic
//	fmt.Println("Solving questions feature is not yet implemented.")
//}
//
//func (u *UI) ViewDashboard() {
//	// Placeholder for viewing the dashboard logic
//	fmt.Println("Viewing dashboard feature is not yet implemented.")
//}
