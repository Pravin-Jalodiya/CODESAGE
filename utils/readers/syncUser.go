package readers

import (
	"cli-project/config"
	"os"
)

func SyncUser() {
	UserStore = FReaderUser(config.USER_FILE, os.O_RDONLY|os.O_CREATE)
}
