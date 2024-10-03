package ui

import (
	"cli-project/pkg/utils"
	"fmt"
	"strings"
)

func (ui *UI) UpdateProgressPage() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println(utils.Colorize("           UPDATE PROGRESS          ", "cyan", "bold"))
		fmt.Println(utils.Colorize("====================================", "cyan", "bold"))
		fmt.Println(utils.Colorize("1. Update progress", "", ""))
		fmt.Println(utils.Colorize("2. Go back", "", ""))

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
			ui.updateProgress()
		case "2":
			return
		default:
			fmt.Println(utils.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		}

	}
}

func (ui *UI) updateProgress() {

	// Update the user's progress by marking the selected question as done
	fmt.Println(utils.Colorize("Fetching progress updates...", "green", ""))
	err := ui.userService.UpdateUserProgress()

	if err != nil {
		fmt.Println(utils.Colorize("Failed to update progress: ", "red", "bold"), err)
		return
	} else {
		fmt.Println(utils.Colorize("Updated progress successfully", "green", "bold"))
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
