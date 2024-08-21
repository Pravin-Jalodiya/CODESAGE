package globals

import "cli-project/internal/domain/models"

var UserStore = make(map[string]models.User)         // username : User
var QuestionStore = make(map[string]models.Question) // question :
var ActiveUser string
