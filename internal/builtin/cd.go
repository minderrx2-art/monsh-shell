package builtin

import (
	"fmt"
	"os"
)

func Cd(dir string) error {
	err := os.Chdir(dir)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("cd: %s: No such file or directory", dir)
	}
	return nil
}
