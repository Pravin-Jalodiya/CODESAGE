package ui

import (
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils/emojis"
	"cli-project/pkg/utils/formatting"
	"cli-project/pkg/validation"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"time"
)

func (ui *UI) ShowSignupPage() {
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))
	fmt.Println(formatting.Colorize("               SIGNUP                  ", "magenta", "bold"))
	fmt.Println(formatting.Colorize("====================================", "magenta", "bold"))

	// Read Username
	var username string
	for {
		fmt.Print(formatting.Colorize("Username: ", "blue", ""))
		username, _ = ui.reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if validation.ValidateUsername(username) {
			break
		}
		fmt.Println(emojis.Error, "Invalid username. It should be between 4 and 20 characters long, should not be only numbers and contain no spaces.")
	}

	// Read Password
	var password, confirmPassword string
	for {
		fmt.Print(formatting.Colorize("Password: ", "blue", ""))
		passwordBytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password = string(passwordBytes)
		fmt.Println()

		// Read Confirm Password
		fmt.Print(formatting.Colorize("Confirm Password: ", "blue", ""))
		confirmPasswordBytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		confirmPassword = string(confirmPasswordBytes)
		fmt.Println()

		// Validate Passwords
		if password == confirmPassword && validation.ValidatePassword(password) {
			break
		}
		if password != confirmPassword {
			fmt.Println(emojis.Error, "Passwords do not match. Please try again.")
		} else {
			fmt.Println(emojis.Error, "Invalid password. It must be at least 8 characters long and include upper/lowercase letters, a digit, and a special character.")
		}
	}

	// Read Name
	var name string
	for {
		fmt.Print(formatting.Colorize("Name: ", "blue", ""))
		name, _ = ui.reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if validation.ValidateName(name) {
			break
		}
		fmt.Println(emojis.Error, "Invalid name. It should be up to 45 characters long and contain only letters and spaces.")
	}

	// Read Email
	var email string
	for {
		fmt.Print(formatting.Colorize("Email: ", "blue", ""))
		email, _ = ui.reader.ReadString('\n')
		email = strings.TrimSpace(email)

		// Validate Email
		if validation.ValidateEmail(email) {
			break
		}
		fmt.Println(emojis.Error, "Invalid email format.")
	}

	// Read LeetCode Username
	var leetcodeID string
	for {
		fmt.Print(formatting.Colorize("LeetCode Username: ", "blue", ""))
		leetcodeID, _ = ui.reader.ReadString('\n')
		leetcodeID = strings.TrimSpace(leetcodeID)

		// Validate LeetCode Username
		exists, err := validation.ValidateLeetCodeUsername(leetcodeID)
		if err != nil {
			fmt.Println(emojis.Error, "Error validating LeetCode username:", err)
			continue
		}
		if exists {
			break
		}
		fmt.Println(emojis.Error, "LeetCode username does not exist.")
	}

	// Create User Object
	user := models.StandardUser{
		StandardUser: models.User{
			Username: username,
			Password: password,
			Name:     name,
			Email:    email,
		},
		LeetcodeID:      leetcodeID,
		QuestionsSolved: []int{},
		LastSeen:        time.Now().UTC(),
	}

	// Call Signup Service
	err := ui.userService.SignUp(user)
	if err != nil {
		fmt.Println(emojis.Error, "Signup failed:", err)
		return
	}

	fmt.Println(emojis.Success, "Signup successful!")
	fmt.Println()
}
