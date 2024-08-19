package models

type Admin struct {
	AdminID  int
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role,omitempty"`
}
