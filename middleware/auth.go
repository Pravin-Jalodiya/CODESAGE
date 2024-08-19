package middleware

var ActiveUserID int

func Auth(userID int) {
	ActiveUserID = userID
}
