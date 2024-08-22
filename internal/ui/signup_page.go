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
			unique, err := ui.userService.IsUsernameUnique(username)
			if err != nil {
				fmt.Println(emojis.Error, "Error checking username uniqueness:", err)
				continue
			}
			if unique {
				break
			}
			fmt.Println(emojis.Info, "Username already taken. Choose another username.")
		} else {
			fmt.Println(emojis.Error, "Invalid username. It should be between 4 and 20 characters long, should not be only numbers and contain no spaces.")
		}
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
			fmt.Println(emojis.Error, "Invalid password. It must be at least 8 characters long and include at least 1 upper/lowercase letters, 1 digit, and 1 special character.")
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

		if validation.ValidateEmail(email) {
			unique, err := ui.userService.IsEmailUnique(email)
			if err != nil {
				fmt.Println(emojis.Error, "Error checking email uniqueness:", err)
				continue
			}
			if unique {
				break
			}
			fmt.Println(emojis.Info, "Email already registered. Use a different email.")
		} else {
			fmt.Println(emojis.Error, "Invalid email format.")
		}
	}

	// Read LeetCode Username
	var leetcodeID string
	for {
		fmt.Print(formatting.Colorize("LeetCode Username: ", "blue", ""))
		leetcodeID, _ = ui.reader.ReadString('\n')
		leetcodeID = strings.TrimSpace(leetcodeID)

		// Check if LeetCode ID is unique in the database
		isUnique, err := ui.userService.IsLeetcodeIDUnique(leetcodeID)
		if err != nil {
			fmt.Println(emojis.Error, "Error checking LeetCode ID uniqueness:", err)
			continue
		}
		if !isUnique {
			fmt.Println(emojis.Error, "LeetCode ID is already taken. Please choose a different ID.")
			continue
		}

		// Validate LeetCode Username with LeetCode API
		exists, err := validation.ValidateLeetcodeUsername(leetcodeID)
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
