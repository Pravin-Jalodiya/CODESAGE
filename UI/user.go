package UI

import (
	"fmt"
	"github.com/fatih/color"
)

func User() {
	var choice int

	red := color.New(color.FgRed).SprintFunc()
	//green := color.New(color.FgGreen).SprintFunc()
	//yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	courseEmoji := "📚"
	toDoEmoji := "✅"
	dailyStatusEmoji := "📅"
	progressEmoji := "📊"
	logoutEmoji := "🚪"
	errorEmoji := "❌"
	//successEmoji := "✅"

	for {
		fmt.Printf("\n%s%sHOME%s%s\n\n%sPlease select an option:\n1. %s Manage course list\n2. %s Manage ToDo list\n3. %s Daily status\n4. %s View progress\n5. %s Log out\n",
			cyan("======"), cyan(" "), cyan("======"), cyan(" "),
			blue(""), courseEmoji, toDoEmoji, dailyStatusEmoji, progressEmoji, logoutEmoji)

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println(red(errorEmoji), red("Invalid input:"), err)
			continue
		}

		switch choice {
		case 1:
			//course.Main(middleware.ActiveUser)

		case 2:
			//generalToDo.Main(middleware.ActiveUser)

		case 3:
			//dailyStatus.Main(middleware.ActiveUser)

		case 4:
			//progress.View(middleware.ActiveUser)

		case 5:
			//middleware.ActiveUser = ""
			fmt.Println(blue(logoutEmoji), blue("User Logged out"))
			return

		default:
			fmt.Println(red(errorEmoji), red("Invalid selection. Please try again."))
		}
	}
}
