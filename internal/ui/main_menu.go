package ui

import (
	"cli-project/pkg/utils"
	"fmt"
	"strings"
)

// ShowMainMenu displays the main menu and handles user input
func (ui *UI) ShowMainMenu() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		// Print the application name and menu options
		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println(utils.Colorize("              CODESAGE              ", "cyan", "bold"))
		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println("Please choose an option:")
		fmt.Printf("1. %s Sign Up\n", utils.SignupEmoji)
		fmt.Printf("2. %s Login\n", utils.LoginEmoji)
		fmt.Printf("3. %s Exit\n", utils.ExitEmoji)
		fmt.Print("Enter your choice : ")

		// Read user input
		choice, _ := ui.reader.ReadString('\n')
		choice = strings.TrimSuffix(choice, "\n")
		choice = strings.TrimSpace(choice)
		fmt.Println()

		switch choice {
		case "1":
			ui.ShowSignupPage()
		case "2":
			ui.ShowLoginPage()
		case "3":
			fmt.Println(utils.ExitEmoji + " Exiting the application.")
			return
		default:
			fmt.Println(utils.ErrorEmoji + " Invalid choice. Please try again.")
		}
	}
}
