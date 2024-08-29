package interfaces

type AuthService interface {
	IsEmailUnique(email string) (bool, error)
	IsUsernameUnique(username string) (bool, error)
	IsLeetcodeIDUnique(leetcodeID string) (bool, error)
	ValidateLeetcodeUsername(username string) (bool, error)
}
