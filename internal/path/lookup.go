package path

import (
	"fmt"
	"os/exec"
)

func Find(binary string) (bool, string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return false, "", fmt.Errorf("Not found")
	}
	return true, path, nil
}
