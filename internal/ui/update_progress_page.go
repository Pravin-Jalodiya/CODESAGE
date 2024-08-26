package ui

import (
	"cli-project/pkg/utils/data_cleaning"
	"cli-project/pkg/utils/formatting"
	"cli-project/pkg/validation"
	"fmt"
	"strings"
)

func (ui *UI) UpdateProgressPage() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
		fmt.Println(formatting.Colorize("           UPDATE PROGRESS          ", "magenta", "bold"))
		fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
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

	var questionID string
	var err error

	for {
		fmt.Print("Enter the ID of the question: ")
		questionID, err = ui.reader.ReadString('\n')
		questionID = strings.TrimSuffix(questionID, "\n")
		questionID = data_cleaning.CleanString(questionID)
		valid, err := validation.ValidateQuestionID(questionID)
		if !valid {
			fmt.Println(err)
			continue
		}
		break
	}
	// Update the user's progress by marking the selected question as done
	progressUpdated, err := ui.userService.UpdateUserProgress(questionID)

	if err != nil {
		fmt.Println(formatting.Colorize("Failed to update progress: ", "red", "bold"), err)
		return
	} else if !progressUpdated {
		fmt.Println(formatting.Colorize("Question already marked as done", "yellow", "bold"))
	} else {
		fmt.Println(formatting.Colorize("Updated progress successfully", "green", "bold"))
	}
}
