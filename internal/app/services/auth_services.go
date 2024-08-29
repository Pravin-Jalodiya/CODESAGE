package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/domain/interfaces"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	leetcodeAPI interfaces2.LeetcodeAPI
}

func NewAuthService(userRepo interfaces.UserRepository, leetcodeAPI interfaces2.LeetcodeAPI) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		leetcodeAPI: leetcodeAPI,
	}
}

func (s *AuthService) IsEmailUnique(email string) (bool, error) {
	return s.userRepo.IsEmailUnique(email)
}

func (s *AuthService) IsUsernameUnique(username string) (bool, error) {
	return s.userRepo.IsUsernameUnique(username)
}

func (s *AuthService) IsLeetcodeIDUnique(leetcodeID string) (bool, error) {
	return s.userRepo.IsLeetcodeIDUnique(leetcodeID)
}

// ValidateLeetcodeUsername checks if the provided LeetCode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	return s.leetcodeAPI.ValidateUsername(username)
}
