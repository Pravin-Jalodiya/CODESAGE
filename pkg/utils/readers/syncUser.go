package readers

import (
	"cli-project/internal/config"
	"cli-project/pkg/globals"
	"os"
)

func SyncUserStore() {
	globals.UserStore = FReaderUser(config.USER_FILE, os.O_RDONLY|os.O_CREATE)
}
