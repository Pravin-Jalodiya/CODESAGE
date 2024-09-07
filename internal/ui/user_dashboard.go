package ui

import (
	"cli-project/pkg/globals"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

// ShowUserDashboard displays the user's Leetcode stats and Codesage stats on the dashboard.
func (ui *UI) ShowUserDashboard() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))
	fmt.Println(formatting.Colorize("            USER DASHBOARD          ", "cyan", "bold"))
	fmt.Println(formatting.Colorize("====================================", "cyan", "bold"))

	// Fetch Leetcode stats
	leetcodeStats, err := ui.userService.GetUserLeetcodeStats(globals.ActiveUserID)
	if err != nil {
		fmt.Println("Error fetching leetcodeStats:", err)
		return
	}

	// Fetch Codesage stats
	codesageStats, err := ui.userService.GetUserCodesageStats(globals.ActiveUserID)
	if err != nil {
		fmt.Println("Error fetching Codesage stats:", err)
		return
	}

	// Create a new table for Leetcode stats
	fmt.Println(formatting.Colorize("Leetcode Stats", "cyan", "bold"))
	leetcodeTable := tablewriter.NewWriter(os.Stdout)
	leetcodeTable.SetHeader([]string{"Difficulty", "Solved", "Total"})

	// Apply color coding to rows
	leetcodeTable.Rich([]string{
		"Easy", fmt.Sprintf("%d", leetcodeStats.EasyDoneCount), fmt.Sprintf("%d", leetcodeStats.TotalEasyCount),
	}, []tablewriter.Colors{{tablewriter.FgGreenColor, tablewriter.Bold}, {tablewriter.FgGreenColor}, {tablewriter.FgGreenColor}})

	leetcodeTable.Rich([]string{
		"Medium", fmt.Sprintf("%d", leetcodeStats.MediumDoneCount), fmt.Sprintf("%d", leetcodeStats.TotalMediumCount),
	}, []tablewriter.Colors{{tablewriter.FgYellowColor, tablewriter.Bold}, {tablewriter.FgYellowColor}, {tablewriter.FgYellowColor}})

	leetcodeTable.Rich([]string{
		"Hard", fmt.Sprintf("%d", leetcodeStats.HardDoneCount), fmt.Sprintf("%d", leetcodeStats.TotalHardCount),
	}, []tablewriter.Colors{{tablewriter.FgRedColor, tablewriter.Bold}, {tablewriter.FgRedColor}, {tablewriter.FgRedColor}})

	leetcodeTable.Render()

	// Create a table for recent accepted submissions
	fmt.Println(formatting.Colorize("\nRecent Accepted Submissions", "cyan", "bold"))
	submissionsTable := tablewriter.NewWriter(os.Stdout)
	submissionsTable.SetHeader([]string{"Recent AC Submissions"})
	for _, submission := range leetcodeStats.RecentACSubmissionTitles {
		submissionsTable.Append([]string{submission})
	}
	submissionsTable.Render()

	// Create a new table for Codesage stats
	fmt.Println(formatting.Colorize("\nCodesage Stats", "cyan", "bold"))
	codesageTable := tablewriter.NewWriter(os.Stdout)
	codesageTable.SetHeader([]string{"Metric", "Count"})

	// Append rows with color coding
	codesageTable.Rich([]string{
		"Total Questions on Platform", fmt.Sprintf("%d", codesageStats.TotalQuestionsCount),
	}, []tablewriter.Colors{{tablewriter.FgMagentaColor}, {tablewriter.FgCyanColor}})

	codesageTable.Rich([]string{
		"Total Questions Solved", fmt.Sprintf("%d", codesageStats.TotalQuestionsDoneCount),
	}, []tablewriter.Colors{{tablewriter.FgMagentaColor}, {tablewriter.FgCyanColor}})

	codesageTable.Rich([]string{
		"Easy Solved", fmt.Sprintf("%d", codesageStats.EasyDoneCount),
	}, []tablewriter.Colors{{tablewriter.FgGreenColor, tablewriter.Bold}, {tablewriter.FgGreenColor}})

	codesageTable.Rich([]string{
		"Medium Solved", fmt.Sprintf("%d", codesageStats.MediumDoneCount),
	}, []tablewriter.Colors{{tablewriter.FgYellowColor, tablewriter.Bold}, {tablewriter.FgYellowColor}})

	codesageTable.Rich([]string{
		"Hard Solved", fmt.Sprintf("%d", codesageStats.HardDoneCount),
	}, []tablewriter.Colors{{tablewriter.FgRedColor, tablewriter.Bold}, {tablewriter.FgRedColor}})

	codesageTable.Render()

	// Display topic-wise stats in a table
	fmt.Println(formatting.Colorize("\nTopic-wise Stats", "cyan", "bold"))
	topicTable := tablewriter.NewWriter(os.Stdout)
	topicTable.SetHeader([]string{"Topic", "Solved"})
	for topic, count := range codesageStats.TopicWiseStats {
		topicTable.Append([]string{topic, fmt.Sprintf("%d", count)})
	}
	topicTable.Render()

	// Display company-wise stats in a table
	fmt.Println(formatting.Colorize("\nCompany-wise Stats", "cyan", "bold"))
	companyTable := tablewriter.NewWriter(os.Stdout)
	companyTable.SetHeader([]string{"Company", "Solved"})
	for company, count := range codesageStats.CompanyWiseStats {
		companyTable.Append([]string{company, fmt.Sprintf("%d", count)})
	}
	companyTable.Render()

	fmt.Println("\nPress any key to go back...")
	_, _ = ui.reader.ReadString('\n')
}
