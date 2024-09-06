package ui

import (
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strings"
)

func (ui *UI) UpdateProgressPage() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("           UPDATE PROGRESS          ", "cyan", "bold"))
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("1. Update progress", "", ""))
		fmt.Println(formatting.Colorize("2. Go back", "", ""))

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
			ui.updateProgress()
		case "2":
			return
		default:
			fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		}

	}
}

func (ui *UI) updateProgress() {

	// Update the user's progress by marking the selected question as done
	fmt.Println(formatting.Colorize("Fetching progress updates...", "green", ""))
	err := ui.userService.UpdateUserProgress()

	if err != nil {
		fmt.Println(formatting.Colorize("Failed to update progress: ", "red", "bold"), err)
		return
	} else {
		fmt.Println(formatting.Colorize("Updated progress successfully", "green", "bold"))
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
