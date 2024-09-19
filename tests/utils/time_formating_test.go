package utils

import (
	"cli-project/pkg/utils"
	"testing"
	"time"
)

func TestConvertToIST(t *testing.T) {

	// Test case: converting from UTC to IST
	utcTime, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z") // This is 00:00:00 on 01/01/2000 in UTC.
	expectedIST := "01/01/2000 05:30:00"                           // Corresponding IST time.

	result := utils.ConvertToIST(utcTime)
	if result != expectedIST {
		t.Errorf("Expected %s, but got %s", expectedIST, result)
	}

	est, _ := time.LoadLocation("America/New_York")
	estTime := time.Date(2000, 1, 1, 0, 0, 0, 0, est) // 01/01/2000 00:00:00 EST
	expectedIST = "01/01/2000 10:30:00"               // Corresponding IST time.

	result = utils.ConvertToIST(estTime)
	if result != expectedIST {
		t.Errorf("Expected %s, but got %s", expectedIST, result)
	}
}
