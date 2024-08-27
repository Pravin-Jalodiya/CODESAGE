package ui

import (
	"bufio"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"
)

func (ui *UI) ShowAdminDashboard() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("           ADMIN DASHBOARD          ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))

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

	// Display the counts
	fmt.Println(formatting.Colorize(fmt.Sprintf("Active Users (Last 24 Hours): %d", activeUsers), "cyan", "bold"))
	fmt.Println(formatting.Colorize(fmt.Sprintf("Total Questions on the Platform: %d", totalQuestions), "cyan", "bold"))

	fmt.Println("\nPress any key to go back...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
