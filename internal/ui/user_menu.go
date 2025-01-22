package ui

import (
	"cli-project/pkg/utils"
	"fmt"
	"strings"
)

func (ui *UI) ShowUserMenu() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println(utils.Colorize("              USER MENU             ", "cyan", "bold"))
		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println(utils.Colorize("1. View questions", "", ""))
		fmt.Println(utils.Colorize("2. View dashboard", "", ""))
		fmt.Println(utils.Colorize("3. Update progress", "", ""))
		fmt.Println(utils.Colorize("4. View profile", "", ""))
		fmt.Println(utils.Colorize("5. Logout", "", ""))

		fmt.Print(utils.Colorize("Enter your choice: ", "yellow", "bold"))
		choice, err := ui.reader.ReadString('\n')
		choice = strings.TrimSuffix(choice, "\n")
		choice = strings.TrimSpace(choice)
		if err != nil {
			fmt.Println(utils.Colorize("Error reading input:", "red", "bold"), err)
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
			ui.ShowUserProfile()
		case "5":
			err := ui.userService.Logout()
			if err != nil {
				fmt.Println(utils.Colorize("Error logging out: ", "red", "bold"), err)
			} else {
				fmt.Printf("%s Logging out...\n", utils.InfoEmoji)
				return
			}
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
