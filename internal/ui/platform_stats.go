package ui

import (
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func (ui *UI) DisplayPlatformStats() {
	fmt.Print("\033[H\033[2J")

	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("        📊 PLATFORM STATS 📊        ", "cyan", "bold"))
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))

	// Fetch platform stats
	platformStats, err := ui.userService.GetPlatformStats()
	if err != nil {
		fmt.Println(formatting.Colorize("Error fetching platform stats: ", "red", "bold"), err)
		return
	}

	// Create and display table for active users and total questions
	statsTable := tablewriter.NewWriter(os.Stdout)
	statsTable.SetHeader([]string{"Metric", "Value"})
	statsTable.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})
	statsTable.SetColWidth(40) // Set a fixed column width

	statsTable.Append([]string{
		"Active Users (Last 24 Hours)",
		formatting.Colorize(fmt.Sprintf("%d", platformStats.ActiveUserInLast24Hours), "cyan", "bold"),
	})
	statsTable.Append([]string{
		"Total Questions on the Platform",
		formatting.Colorize(fmt.Sprintf("%d", platformStats.TotalQuestionsCount), "cyan", "bold"),
	})

	statsTable.SetAutoWrapText(true)
	statsTable.SetRowLine(true)
	statsTable.Render()

	// Create and display table for difficulty-wise question counts
	fmt.Println(formatting.Colorize("\nDifficulty-wise Questions Count", "cyan", "bold"))
	difficultyTable := tablewriter.NewWriter(os.Stdout)
	difficultyTable.SetHeader([]string{"Difficulty", "Count"})
	difficultyTable.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

	// Define the fixed order of difficulties
	difficultyOrder := []string{"easy", "medium", "hard"}

	for _, difficulty := range difficultyOrder {
		count, exists := platformStats.DifficultyWiseQuestionsCount[difficulty]
		if !exists {
			count = 0 // Default to 0 if the difficulty level is not present in the map
		}

		color := "white" // default color
		switch difficulty {
		case "easy":
			color = "green"
		case "medium":
			color = "yellow"
		case "hard":
			color = "red"
		}

		difficultyTable.Append([]string{
			formatting.Colorize(difficulty, color, "bold"),
			fmt.Sprintf("%d", count),
		})
	}

	difficultyTable.SetAutoWrapText(true)
	difficultyTable.SetRowLine(true)
	difficultyTable.Render()

	// Create and display table for topic-wise question counts
	fmt.Println(formatting.Colorize("\nTopic-wise Questions Count", "cyan", "bold"))
	topicTable := tablewriter.NewWriter(os.Stdout)
	topicTable.SetHeader([]string{"Topic", "Count"})
	topicTable.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

	for topic, count := range platformStats.TopicWiseQuestionsCount {
		topicTable.Append([]string{topic, fmt.Sprintf("%d", count)})
	}

	topicTable.SetAutoWrapText(true)
	topicTable.SetRowLine(true)
	topicTable.Render()

	// Create and display table for company-wise question counts
	fmt.Println(formatting.Colorize("\nCompany-wise Questions Count", "cyan", "bold"))
	companyTable := tablewriter.NewWriter(os.Stdout)
	companyTable.SetHeader([]string{"Company", "Count"})
	companyTable.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

	for company, count := range platformStats.CompanyWiseQuestionsCount {
		companyTable.Append([]string{company, fmt.Sprintf("%d", count)})
	}

	companyTable.SetAutoWrapText(true)
	companyTable.SetRowLine(true)
	companyTable.Render()

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
