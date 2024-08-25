package ui

import (
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func (ui *UI) ViewQuestionsPage() {

	for {

		// Clear the screen
		fmt.Print("\033[H\033[2J")

		fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
		fmt.Println(formatting.Colorize("              QUESTIONS             ", "magenta", "bold"))
		fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
		fmt.Printf("1. %s View questions\n", emojis.View)
		fmt.Printf("2. %s Go back\n", emojis.Back)
		fmt.Print("Enter your choice : ")

		// Read user input
		choice, _ := ui.reader.ReadString('\n')
		choice = strings.TrimSuffix(choice, "\n")
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			ui.ViewQuestions()
		case "2":
			return
		default:
			fmt.Println(formatting.Colorize("Invalid choice. Please select a valid option.", "red", "bold"))
		}

	}

}

func (ui *UI) ViewQuestions() {

	// Load all questions in the db
	questionsList, err := ui.questionService.GetAllQuestions()
	if err != nil {
		fmt.Println("Failed to load questions")
		return
	}

	// If no questions found, notify the user
	if len(questionsList) == 0 {
		fmt.Println("Trouble loading questions. Try again later.")
		return
	}

	// Create a new tab writer to format the output as a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	// Print table headers
	_, err = fmt.Fprintln(w, "ID\tTitle\tDifficulty\tLink\tTopic-Tags\tCompany-Tags")
	if err != nil {
		fmt.Println("Error rendering page.")
	}

	// Print table rows
	for _, question := range questionsList {
		// Convert slices to comma-separated strings for display
		topicTags := fmt.Sprintf("%v", question.TopicTags)
		companyTags := fmt.Sprintf("%v", question.CompanyTags)

		// Format the question details into table rows
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			question.QuestionID,
			question.QuestionTitle,
			question.Difficulty,
			question.QuestionLink,
			topicTags[1:len(topicTags)-1], // Remove square brackets from slice string
			companyTags[1:len(companyTags)-1],
		)

		if err != nil {
			fmt.Println("Error rendering page.")
			return
		}
	}

	// Flush the writer to ensure all output is printed
	err = w.Flush()
	if err != nil {
		fmt.Println("Error rendering page.")
		return
	}
}

func (ui *UI) ViewFilteredQuestions(difficulty, topicTag, companyTag string) {

}
