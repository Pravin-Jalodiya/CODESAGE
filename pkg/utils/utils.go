package utils

import (
	"cli-project/internal/config"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	SignupEmoji   = "✍️"
	LoginEmoji    = "🔑"
	ExitEmoji     = "🚪"
	ErrorEmoji    = "❌"
	SuccessEmoji  = "✅"
	ProfileEmoji  = "👤"
	StatsEmoji    = "📊"
	SettingsEmoji = "⚙️"
	QuestionEmoji = "❓"
	InfoEmoji     = "ℹ️"
	BackEmoji     = "🔙"
	ViewEmoji     = "👁️"
)

// Colorize generates a colorized string with the specified foreground color and style
func Colorize(text, fgColor, style string) string {
	var c *color.Color

	// Apply foreground color
	switch fgColor {
	case "black":
		c = color.New(color.FgBlack)
	case "red":
		c = color.New(color.FgRed)
	case "green":
		c = color.New(color.FgGreen)
	case "yellow":
		c = color.New(color.FgYellow)
	case "blue":
		c = color.New(color.FgBlue)
	case "magenta":
		c = color.New(color.FgMagenta)
	case "cyan":
		c = color.New(color.FgCyan)
	case "white":
		c = color.New(color.FgWhite)
	default:
		c = color.New(color.Reset)
	}

	// Apply style
	switch style {
	case "bold":
		c = c.Add(color.Bold)
	case "underline":
		c = c.Add(color.Underline)
	}

	// Return the formatted string
	return c.Sprint(text)
}

func AreSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func ReadCSV(filePath string) ([][]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening question file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file")
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CapitalizeWords capitalizes the first letter of each word in a string.
func CapitalizeWords(s string) string {
	words := strings.Fields(s) // Split the string into words
	for i, word := range words {
		// Capitalize the first letter of each word
		if len(word) > 0 {
			words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func CleanString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func CleanTags(tags string) []string {
	tagList := strings.Split(tags, ",")
	for i, tag := range tagList {
		tagList[i] = CleanString(tag)
	}
	return tagList
}

// HashString generates a bcrypt hash for the given password.
func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyString verifies if the given password matches the stored hash.
func VerifyString(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateJwtToken generates a JWT token for the given user ID and role
func CreateJwtToken(username string, userId string, role string) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"role":     role,
		"exp":      time.Now().Add(5 * time.Minute).Unix(), // Token expiry time (1 minute)
	}

	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(config.SECRET_KEY)
	if err != nil {
		// Log the error if needed (uncomment the line below)
		// logger.Logger.Errorw("Error signing token", "error", err, "time", time.Now())
		return "", errors.New("error creating jwt token")
	}

	return tokenString, nil
}

func MergeTags(existingTags, newTags []string) []string {
	tagSet := make(map[string]struct{})

	for _, tag := range existingTags {
		tagSet[tag] = struct{}{}
	}

	for _, tag := range newTags {
		if _, found := tagSet[tag]; !found {
			existingTags = append(existingTags, tag)
		}
	}

	return existingTags
}

// ConvertToIST converts a UTC time to IST and returns it in dd/mm/yyyy hh:mm:ss format.
func ConvertToIST(t time.Time) string {
	// Define IST timezone
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return "Invalid time location"
	}

	// Convert UTC time to IST
	istTime := t.In(istLocation)

	// Format the time as dd/mm/yyyy hh:mm:ss
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d",
		istTime.Day(),
		istTime.Month(),
		istTime.Year(),
		istTime.Hour(),
		istTime.Minute(),
		istTime.Second())
}

func GenerateUUID() string {
	return uuid.New().String()
}
