package services

import "cli-project/internal/domain/interfaces"

type AuthService struct {
	userRepo interfaces.UserRepository
}

func NewAuthService(userRepo interfaces.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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
