package ui

import (
	"cli-project/internal/config"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (ui *UI) ManageQuestions() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("          MANAGE QUESTIONS          ", "cyan", "bold"))
		fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
		fmt.Println(formatting.Colorize("1. Add questions", "", ""))
		fmt.Println(formatting.Colorize("2. Remove question", "", ""))
		fmt.Println(formatting.Colorize("3. Go back", "", ""))

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
			ui.AddQuestions()
		case "2":
			ui.RemoveQuestion()
		case "3":
			return
		default:
			fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		}
	}
}

func (ui *UI) AddQuestions() {
	// List all files in the CSV_DIR
	files, err := os.ReadDir(config.CSV_DIR)
	if err != nil {
		fmt.Println(formatting.Colorize("Error reading directory", "red", "bold"))
		return
	}

	if len(files) == 0 {
		fmt.Println(formatting.Colorize("No files found in the directory.", "yellow", "bold"))
		return
	}

	fmt.Println(formatting.Colorize("Select a file to add questions from:", "cyan", "bold"))

	// Display the list of files with their extensions
	for i, file := range files {
		fmt.Printf("[%d] %s\n", i+1, file.Name())
	}

	// Ask admin to select a file
	var choice int
	var fileName string
	for {

		fmt.Print(formatting.Colorize("Enter the file number: ", "yellow", "bold"))
		_, err = fmt.Scan(&choice)
		if err != nil || choice < 1 || choice > len(files) {
			fmt.Println(formatting.Colorize("Invalid choice.", "red", "bold"))
			continue
		}

		selectedFile := files[choice-1]
		fileName = selectedFile.Name()

		// Check if the file has a .csv extension
		if filepath.Ext(fileName) != ".csv" {
			fmt.Println(formatting.Colorize("Selected file is not a CSV file.", "red", "bold"))
			continue
		}
		break
	}

	// Construct the full path to the selected file
	fullFilePath := filepath.Join(config.CSV_DIR, fileName)

	// Call the service method to add questions from the selected file
	newQuestionsAdded, err := ui.questionService.AddQuestionsFromFile(fullFilePath)
	if err != nil {
		fmt.Println(formatting.Colorize("Error adding questions from file:", "red", "bold"), fileName, err)
		return
	} else if !newQuestionsAdded {
		fmt.Println(formatting.Colorize("No new questions in the file:", "yellow", "bold"), fileName)
	} else {
		fmt.Println(formatting.Colorize("Questions successfully added from file:", "green", "bold"), fileName)
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')

}

func (ui *UI) RemoveQuestion() {

	ui.ViewQuestions()
	// Placeholder for the remove question logic
	// Prompt admin to enter the Question ID
	var questionID string
	var err error
	for {
		fmt.Print(formatting.Colorize("Enter the Question ID to remove: ", "yellow", "bold"))
		questionID, err = ui.reader.ReadString('\n')
		questionID = strings.TrimSuffix(questionID, "\n")
		questionID = strings.TrimSpace(questionID)
		if err != nil {
			fmt.Println(formatting.Colorize("Error reading input.", "red", "bold"))
			continue
		}
		break
	}
	// Call the QuestionService to remove the question
	err = ui.questionService.RemoveQuestionByID(questionID)
	if err != nil {
		fmt.Println(formatting.Colorize("Failed to remove the question:", "red", "bold"), err)
		return
	} else {
		fmt.Println(formatting.Colorize("Question removed successfully!", "green", "bold"))
	}

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
