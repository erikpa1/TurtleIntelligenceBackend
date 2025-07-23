package credentials

import (
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"os/user"
	"path/filepath"
)

func GetDarwinWorkspace() string {

	fromEnv := GetEnvOrDefault(fmt.Sprintf("TURTLE_%s_MASTER", APP_NAME), "")

	if fromEnv == "" {
		usr, err := user.Current()
		appSupportPath := filepath.Join(usr.HomeDir, "Library", "Application Support")
		if err == nil {
			return appSupportPath + "/" + GetAppName()
		} else {
			lg.LogE(err)
		}

	} else {
		return fromEnv + "/" + GetAppName()
	}

	return "../TurtleIntelligenceStorage" + "/" + GetAppName()

}

func GetStoragePath() string {
	return GetEnvOrDefault(fmt.Sprintf("TURTLE_%s_STORAGE", APP_NAME), "LOCALAPPDATA")
}
