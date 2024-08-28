package ui

import (
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strings"
)

// ShowMainMenu displays the main menu and handles user input
func (ui *UI) ShowMainMenu() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		// Print the application name and menu options
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("              CODESAGE              ", "cyan", "bold"))
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println("Please choose an option:")
		fmt.Printf("1. %s Sign Up\n", emojis.SignUp)
		fmt.Printf("2. %s Login\n", emojis.Login)
		fmt.Printf("3. %s Exit\n", emojis.Exit)
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
			fmt.Println(emojis.Exit + " Exiting the application.")
			return
		default:
			fmt.Println(emojis.Error + " Invalid choice. Please try again.")
		}
	}
}
