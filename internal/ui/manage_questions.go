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
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("         MANAGE QUESTIONS           ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("1. Add Questions", "green", ""))
	fmt.Println(formatting.Colorize("2. Remove Question", "green", ""))

	fmt.Print(formatting.Colorize("Enter your choice: ", "yellow", "bold"))
	choice, err := ui.reader.ReadString('\n')
	if err != nil {
		fmt.Println(formatting.Colorize("Error reading input:", "red", "bold"), err)
		return
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		ui.AddQuestions()
	case "2":
		ui.RemoveQuestion()
	default:
		fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		ui.ManageQuestions()
	}
	return
}

func (ui *UI) AddQuestions() {
	// List all files in the CSV_DIR
	files, err := os.ReadDir(config.CSV_DIR)
	if err != nil {
		fmt.Println(formatting.Colorize("Error reading directory:", "red", "bold"), err)
		return
	}

	if len(files) == 0 {
		fmt.Println(formatting.Colorize("No files found in the directory.", "yellow", "bold"))
		return
	}

	fmt.Println(formatting.Colorize("Select a file to add questions from:", "green", "bold"))

	// Display the list of files with their extensions
	for i, file := range files {
		fmt.Printf("[%d] %s\n", i+1, file.Name())
	}

	// Ask admin to select a file
	fmt.Print(formatting.Colorize("Enter the file number: ", "blue", "bold"))
	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil || choice < 1 || choice > len(files) {
		fmt.Println(formatting.Colorize("Invalid choice.", "red", "bold"))
		return
	}

	selectedFile := files[choice-1]
	fileName := selectedFile.Name()

	// Check if the file has a .csv extension
	if filepath.Ext(fileName) != ".csv" {
		fmt.Println(formatting.Colorize("Selected file is not a CSV file.", "red", "bold"))
		return
	}

	// Construct the full path to the selected file
	fullFilePath := filepath.Join(config.CSV_DIR, fileName)

	// Call the service method to add questions from the selected file
	err = ui.questionService.AddQuestionsFromFile(fullFilePath)
	if err != nil {
		fmt.Println(formatting.Colorize("Error adding questions from file:", "red", "bold"), err)
	} else {
		fmt.Println(formatting.Colorize("Questions successfully added from file:", "green", "bold"), fileName)
	}
}

func (ui *UI) RemoveQuestion() {
	// Placeholder for the remove question logic
	// Prompt admin to enter the Question ID
	var questionID string
	var err error
	for {
		fmt.Print(formatting.Colorize("Enter the Question ID to remove: ", "yellow", "bold"))
		questionID, err = ui.reader.ReadString('\n')
		if err != nil {
			fmt.Println(formatting.Colorize("Error reading input:", "red", "bold"), err)
			return
		}

		break
	}
	// Call the QuestionService to remove the question
	err = ui.questionService.RemoveQuestionByID(questionID)
	if err != nil {
		fmt.Println(formatting.Colorize("Failed to remove the question:", "red", "bold"), err)
	} else {
		fmt.Println(formatting.Colorize("Question removed successfully!", "green", "bold"))
	}
}
