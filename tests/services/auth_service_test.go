package service_test

import (
	"testing"
)

func TestIsEmailUnique(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	email := "test@example.com"

	// Set up expectations
	mockUserRepo.EXPECT().
		IsEmailUnique(email).
		Return(true, nil).
		Times(1)

	// Call the method
	isUnique, err := authService.IsEmailUnique(email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isUnique {
		t.Fatalf("expected email to be unique")
	}
}

func TestIsUsernameUnique(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "testuser"

	// Set up expectations
	mockUserRepo.EXPECT().
		IsUsernameUnique(username).
		Return(true, nil).
		Times(1)

	// Call the method
	isUnique, err := authService.IsUsernameUnique(username)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isUnique {
		t.Fatalf("expected username to be unique")
	}
}

func TestIsLeetcodeIDUnique(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	LeetcodeID := "Leetcode123"

	// Set up expectations
	mockUserRepo.EXPECT().
		IsLeetcodeIDUnique(LeetcodeID).
		Return(true, nil).
		Times(1)

	// Call the method
	isUnique, err := authService.IsLeetcodeIDUnique(LeetcodeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isUnique {
		t.Fatalf("expected Leetcode ID to be unique")
	}
}

func TestValidateLeetcodeUsername(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "validuser"

	// Set up expectations
	mockLeetcodeAPI.EXPECT().
		ValidateUsername(username).
		Return(true, nil).
		Times(1)

	// Call the method
	isValid, err := authService.ValidateLeetcodeUsername(username)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isValid {
		t.Fatalf("expected Leetcode username to be valid")
	}
}
