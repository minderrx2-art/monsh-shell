package builtin

import (
	"fmt"
	"os"
)

func Cd(dir string) error {
	if dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", dir)
		}
		os.Chdir(homeDir)
	} else {
		err := os.Chdir(dir)
		if err != nil && os.IsNotExist(err) {
			return fmt.Errorf("cd: %s: No such file or directory", dir)
		}
	}
	return nil
}
