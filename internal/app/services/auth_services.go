package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/domain/interfaces"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	LeetcodeAPI interfaces2.LeetcodeAPI
}

func NewAuthService(userRepo interfaces.UserRepository, LeetcodeAPI interfaces2.LeetcodeAPI) interfaces.AuthService {
	return &AuthService{
		userRepo:    userRepo,
		LeetcodeAPI: LeetcodeAPI,
	}
}

func (s *AuthService) IsEmailUnique(email string) (bool, error) {
	return s.userRepo.IsEmailUnique(email)
}

func (s *AuthService) IsUsernameUnique(username string) (bool, error) {
	return s.userRepo.IsUsernameUnique(username)
}

func (s *AuthService) IsLeetcodeIDUnique(LeetcodeID string) (bool, error) {
	return s.userRepo.IsLeetcodeIDUnique(LeetcodeID)
}

// ValidateLeetcodeUsername checks if the provided Leetcode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	return s.LeetcodeAPI.ValidateUsername(username)
}
