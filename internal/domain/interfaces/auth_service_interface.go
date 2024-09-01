package interfaces

type AuthService interface {
	IsEmailUnique(email string) (bool, error)
	IsUsernameUnique(username string) (bool, error)
	IsLeetcodeIDUnique(LeetcodeID string) (bool, error)
	ValidateLeetcodeUsername(username string) (bool, error)
}
