package ui

import (
	"cli-project/pkg/utils"
	"fmt"
)

func (ui *UI) ShowBannedMessage() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Display the banned message
	fmt.Println(utils.Colorize("====================================", "red", "bold"))
	fmt.Println(utils.Colorize("         YOU ARE BANNED FROM        ", "red", "bold"))
	fmt.Println(utils.Colorize("            THE PLATFORM            ", "red", "bold"))
	fmt.Println(utils.Colorize("====================================", "red", "bold"))

	fmt.Println("\nPress any key to go back...")

	_, _ = ui.reader.ReadString('\n')
}
