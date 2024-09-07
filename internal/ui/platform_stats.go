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

	// Fetch the number of active users in the last 24 hours
	activeUsers, err := ui.userService.CountActiveUserInLast24Hours()
	if err != nil {
		fmt.Println(formatting.Colorize("Error fetching active users count: ", "red", "bold"), err)
		return
	}

	// Fetch the total number of questions on the platform
	totalQuestions, err := ui.questionService.GetTotalQuestionsCount()
	if err != nil {
		fmt.Println(formatting.Colorize("Error fetching total questions count: ", "red", "bold"), err)
		return
	}

	// Create a table for the stats
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value"})

	// Set column alignment and width
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})
	table.SetColWidth(40) // Set a fixed column width

	// Add the rows to the table
	table.Append([]string{
		"Active Users (Last 24 Hours)",
		formatting.Colorize(fmt.Sprintf("%d", activeUsers), "cyan", "bold"),
	})
	table.Append([]string{
		"Total Questions on the Platform",
		formatting.Colorize(fmt.Sprintf("%d", totalQuestions), "cyan", "bold"),
	})

	table.Render()

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
