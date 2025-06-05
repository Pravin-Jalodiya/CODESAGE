package roles

import "fmt"

// Role represents a role in the system.
type Role int

const (
	// Define the roles as a typed enum.
	USER Role = iota
	ADMIN
)

// String returns the string representation of the Role.
func (r Role) String() string {
	return [...]string{"user", "admin"}[r]
}

// ParseRole converts a string to a Role.
func ParseRole(roleStr string) (Role, error) {
	switch roleStr {
	case "user":
		return USER, nil
	case "admin":
		return ADMIN, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", roleStr)
	}
}
