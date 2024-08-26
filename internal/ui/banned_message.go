package ui

import (
	"bufio"
	"cli-project/pkg/utils/formatting"
	"fmt"
	"os"
)

func (ui *UI) ShowBannedMessage() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Display the banned message
	fmt.Println(formatting.Colorize("====================================", "red", "bold"))
	fmt.Println(formatting.Colorize("         YOU ARE BANNED FROM        ", "red", "bold"))
	fmt.Println(formatting.Colorize("            THE PLATFORM            ", "red", "bold"))
	fmt.Println(formatting.Colorize("====================================", "red", "bold"))

	// Prompt user to press any key to return
	fmt.Println("\nPress any key to return...")

	// Wait for user input
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
