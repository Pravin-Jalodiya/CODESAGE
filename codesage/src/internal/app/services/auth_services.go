package services

import (
	interfaces2 "codesage/external/domain/interfaces"
	"codesage/internal/domain/interfaces"
	"codesage/internal/domain/models"
	"codesage/pkg/errors"
	"codesage/pkg/globals"

	//"codesage/pkg/globals"
	"errors"

	//"codesage/pkg/globals"
	//"codesage/pkg/logger"
	"codesage/pkg/utils"
	"context"
	//"errors"
	"fmt"
	"strings"
	"time"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	userService interfaces.UserService
	LeetcodeAPI interfaces2.LeetcodeAPI
	emailConfig *utils.EmailConfig
}

func NewAuthService(userRepo interfaces.UserRepository, LeetcodeAPI interfaces2.LeetcodeAPI) interfaces.AuthService {
	return &AuthService{
		userRepo: userRepo,
		//userService: userService,
		LeetcodeAPI: LeetcodeAPI,
		emailConfig: utils.NewEmailConfig(),
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(ctx context.Context, user *models.StandardUser) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.Organisation = utils.CapitalizeWords(user.Organisation)
	user.Country = utils.CapitalizeWords(user.Country)
	user.ID = utils.GenerateUUID()

	// Check if the username is unique
	usernameUnique, err := s.userRepo.IsUsernameUnique(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !usernameUnique {
		return fmt.Errorf("%w: %v", errs.ErrUserNameAlreadyExists, user.Username)
	}

	// Check if the email is unique
	emailUnique, err := s.userRepo.IsEmailUnique(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !emailUnique {
		return fmt.Errorf("%w: %v", errs.ErrEmailAlreadyExists, user.Email)
	}

	// Check if Leetcode ID exists on Leetcode or not
	isLeetcodeUsernameValid, leetcodeErr := s.LeetcodeAPI.ValidateLeetcodeUsername(user.LeetcodeID)
	if leetcodeErr != nil {
		return fmt.Errorf("%w: %v", errs.ErrLeetcodeValidationFailed, leetcodeErr)
	}
	if !isLeetcodeUsernameValid {
		return fmt.Errorf("%w: %v", errs.ErrLeetcodeUsernameInvalid, user.LeetcodeID)
	}

	// Check if the Leetcode ID is unique
	leetcodeIDUnique, err := s.userRepo.IsLeetcodeIDUnique(ctx, user.LeetcodeID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !leetcodeIDUnique {
		return fmt.Errorf("%w: %v", errs.ErrLeetcodeIDAlreadyExists, user.LeetcodeID)
	}

	// Get user's leetcode profile picture
	avatar, _ := s.LeetcodeAPI.GetUserAvatar(user.LeetcodeID)
	if avatar != "" {
		user.Avatar = avatar
	} else {
		user.Avatar = "https://assets.leetcode.com/users/default_avatar.jpg"
	}

	hashedPassword, err := utils.HashString(user.Password)
	if err != nil {
		return fmt.Errorf("%w: could not hash password", errs.ErrInternalServerError)
	}
	user.Password = hashedPassword
	user.Role = "user"
	user.IsBanned = false
	user.QuestionsSolved = []string{}
	user.LastSeen = time.Now().UTC()

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: could not register user", errs.ErrDbError)
	}
	return nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.StandardUser, error) {
	username = utils.CleanString(username)

	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}

	if !utils.VerifyString(password, user.Password) {
		return nil, errs.ErrInvalidPassword
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	user, err := s.userRepo.FetchUserByUsername(ctx, globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}

	user.LastSeen = time.Now().UTC()

	err = s.userRepo.UpdateUserDetails(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}

	globals.ActiveUserID = ""
	return nil
}

func (s *AuthService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	unique, err := s.userRepo.IsEmailUnique(ctx, email)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

func (s *AuthService) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	unique, err := s.userRepo.IsUsernameUnique(ctx, username)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

func (s *AuthService) IsLeetcodeIDUnique(ctx context.Context, LeetcodeID string) (bool, error) {
	unique, err := s.userRepo.IsLeetcodeIDUnique(ctx, LeetcodeID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

// ValidateLeetcodeUsername checks if the provided Leetcode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	valid, err := s.LeetcodeAPI.ValidateLeetcodeUsername(username)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrExternalAPI, err)
	}
	return valid, nil
}

//// GenerateAndSendOtp handles OTP generation and sending
//func (a *AuthService) GenerateAndSendOtp(email string) error {
//	// Input validation
//	if email == "" {
//		return fmt.Errorf("email address is required")
//	}
//
//	// Check if user exists
//	user, err := a.userService.GetUserByEmail(context.TODO(), email)
//	if err != nil {
//		logger.Logger.Errorw("Failed to fetch user", "email", email, "error", err)
//		return fmt.Errorf("failed to process request")
//	}
//
//	if user == nil {
//		logger.Logger.Warnw("No user found with email", "email", email)
//		return fmt.Errorf("invalid email address")
//	}
//
//	// Generate OTP
//	otp, err := utils.GenerateOTP()
//	if err != nil {
//		logger.Logger.Errorw("Failed to generate OTP", "email", email, "error", err)
//		return fmt.Errorf("failed to generate verification code")
//	}
//
//	// Save OTP
//	utils.SaveOTP(email, otp)
//
//	// Send OTP email using new email system
//	if err := a.emailConfig.SendOTPEmail(email, otp); err != nil {
//		logger.Logger.Errorw("Failed to send OTP email",
//			"email", email,
//			"error", err,
//			"userId", user.ID,
//		)
//		return fmt.Errorf("failed to send verification code")
//	}
//
//	logger.Logger.Infow("Successfully sent OTP",
//		"email", email,
//		"userId", user.ID,
//	)
//
//	return nil
//}

//func (a *AuthService) UpdateUserPassword(ctx context.Context, email string, password string) error {
//	hashedPassword, _ := utils.HashString(password)
//	return a.userRepo.UpdateUserPassword(ctx, email, hashedPassword)
//}
