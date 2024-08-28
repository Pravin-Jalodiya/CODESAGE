package ui

import (
	"bufio"
	"cli-project/pkg/utils/data_cleaning"
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
	if len(*questionsList) == 0 {
		fmt.Println("Trouble loading questions. Try again later.")
		return
	}

	// Create a new tab writer to format the output as a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Print table headers
	_, err = fmt.Fprintln(w, "ID\tTitle\tDifficulty\tLink\tTopic-Tags\tCompany-Tags")
	if err != nil {
		fmt.Println("Error rendering page.")
	}

	// Print table rows
	for _, question := range *questionsList {
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

	// Prompt user for input
	fmt.Println("Press 'f' to view filtered questions or any other key to return")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return
	}

	// Check the input and call ViewFilteredQuestions if input is 'f'
	if strings.ToLower(input) == "f" {
		ui.ViewFilteredQuestions()
	}
}

func (ui *UI) ViewFilteredQuestions() {

	// Prompt for difficulty
	fmt.Print("Enter difficulty (or 'any' for no filter): ")
	difficulty, _ := ui.reader.ReadString('\n')
	difficulty = strings.TrimSuffix(difficulty, "\n")
	difficulty = data_cleaning.CleanString(difficulty)

	// Prompt for topic
	fmt.Print("Enter topic (or 'any' for no filter): ")
	topic, _ := ui.reader.ReadString('\n')
	topic = strings.TrimSuffix(topic, "\n")
	topic = data_cleaning.CleanString(topic)

	// Prompt for company
	fmt.Print("Enter company (or 'any' for no filter): ")
	company, _ := ui.reader.ReadString('\n')
	company = strings.TrimSuffix(company, "\n")
	company = data_cleaning.CleanString(company)

	// Fetch filtered questions
	filteredQuestions, err := ui.questionService.GetQuestionsByFilters(difficulty, topic, company)
	if err != nil {
		fmt.Printf("Error fetching filtered questions: %v\n", err)
		return
	}

	// If no questions found, notify the user
	if len(*filteredQuestions) == 0 {
		fmt.Println(formatting.Colorize("no questions match the filter", "yellow", "bold"))
		return
	}

	// Display the questions
	// Create a new tab writer to format the output as a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Print table headers
	_, err = fmt.Fprintln(w, "ID\tTitle\tDifficulty\tLink\tTopic-Tags\tCompany-Tags")
	if err != nil {
		fmt.Println("Error rendering page.")
	}

	// Print table rows
	for _, question := range *filteredQuestions {
		//Convert slices to comma-separated strings for display
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

	fmt.Println("\nPress any key to go back...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
