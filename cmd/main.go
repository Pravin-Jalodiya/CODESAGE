package main

import (
	"cli-project/controllers"
	"cli-project/utils/emojis"
	"cli-project/utils/str"
	"fmt"
	"os"
)

func main() {

	var choice int

	for {
		fmt.Printf("\n%s%sCHEATCODE%s%s\n\n%sPlease select an option:\n1. %s Sign Up\n2. %s Log In\n3. %s Exit\n",
			str.Colorize("======", "cyan", "", "bold"),
			str.Colorize(" ", "cyan", "", ""),
			str.Colorize(" ", "cyan", "", ""),
			str.Colorize("======", "cyan", "", "bold"),
			str.Colorize("", "blue", "", ""),
			str.Colorize(emojis.SignUp, "blue", "", "bold"),
			str.Colorize(emojis.Login, "blue", "", "bold"),
			str.Colorize(emojis.Exit, "blue", "", "bold"),
		)

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println(str.Colorize(emojis.Error, "red", "", "bold"), str.Colorize("Invalid input:", "red", "", "bold"), err)
			continue
		}

		switch choice {
		case 1:
			controllers.SignUp()

		case 2:
			controllers.Login()

		case 3:
			fmt.Println(str.Colorize(emojis.Exit, "blue", "", ""), str.Colorize("Exiting...", "blue", "", ""))
			os.Exit(0)

		default:
			fmt.Println(str.Colorize(emojis.Error, "red", "", ""), str.Colorize("Invalid selection. Please try again.", "red", "", ""))
		}
	}
}
