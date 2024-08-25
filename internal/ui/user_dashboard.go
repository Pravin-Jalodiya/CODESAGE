package ui

import "fmt"

func (ui *UI) ShowUserDashboard() {
	for {
		// Clear the screen
		fmt.Print("\033[H\033[2J")
	}
}
