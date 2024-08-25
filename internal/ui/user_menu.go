package ui

import (
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strings"
)

func (ui *UI) ShowUserMenu() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("              USER MENU             ", "cyan", "bold"))
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("1. View questions", "green", ""))
		fmt.Println(formatting.Colorize("2. View dashboard", "green", ""))
		fmt.Println(formatting.Colorize("3. Update progress", "green", ""))
		fmt.Println(formatting.Colorize("4. Logout", "green", ""))

		fmt.Print(formatting.Colorize("Enter your choice: ", "yellow", "bold"))
		choice, err := ui.reader.ReadString('\n')
		choice = strings.TrimSuffix(choice, "\n")
		choice = strings.TrimSpace(choice)
		if err != nil {
			fmt.Println(formatting.Colorize("Error reading input:", "red", "bold"), err)
			return
		}

		switch choice {
		case "1":
			ui.ViewQuestionsPage()
		case "2":
			ui.ShowUserDashboard()
		case "3":
			ui.UpdateProgressPage()
		case "4":
			ui.userService.Logout()
			fmt.Printf("%s Logging out...\n", emojis.Info)
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
