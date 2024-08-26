package utils

import (
	"fmt"
	"time"
)

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
