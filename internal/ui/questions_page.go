package ui

import (
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"
	"text/tabwriter"
)

func (ui *UI) ViewQuestionsPage() {

	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("              QUESTIONS             ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	// Load all questions in the db
	questionsList, err := ui.questionService.GetAllQuestions()
	if err != nil {
		fmt.Println("Failed to load questions:", err)
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
	fmt.Fprintln(w, "ID\tTitle\tDifficulty\tLink\tTopic-Tags\tCompany-Tags")

	// Print table rows
	for _, question := range questionsList {
		// Convert slices to comma-separated strings for display
		topicTags := fmt.Sprintf("%v", question.TopicTags)
		companyTags := fmt.Sprintf("%v", question.CompanyTags)

		// Format the question details into table rows
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			question.QuestionID,
			question.QuestionTitle,
			question.Difficulty,
			question.QuestionLink,
			topicTags[1:len(topicTags)-1], // Remove square brackets from slice string
			companyTags[1:len(companyTags)-1],
		)
	}

	// Flush the writer to ensure all output is printed
	w.Flush()

}

func (ui *UI) ViewFilteredQuestions(difficulty, topicTag, companyTag string) {

}
