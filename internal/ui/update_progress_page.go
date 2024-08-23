package ui

import (
	"cli-project/pkg/globals"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"strconv"
	"strings"
)

func (ui *UI) UpdateProgressPage() {
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("           UPDATE PROGRESS          ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))

	var questionID string
	var err error
	var ID int
	for {
		fmt.Print("Enter the ID of the question: ")
		questionID, err = ui.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Try again.")
			continue
		}

		questionID = strings.TrimSuffix(questionID, "\n")
		questionID = strings.TrimSpace(questionID)
		ID, err = strconv.Atoi(questionID)
		if err != nil {
			fmt.Println("Invalid question ID. Enter a valid number.", err)
			continue
		}
		break
	}

	// Update the user's progress by marking the selected question as done
	err = ui.userService.UpdateUserProgress(globals.ActiveUser, ID)
	if err != nil {
		fmt.Println("Failed to update progress:", err)
	} else {
		fmt.Println("Progress updated successfully!")
	}
}
