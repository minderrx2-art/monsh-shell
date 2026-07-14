package path

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func scan(prefix string) []string {
	PATH := os.Getenv("PATH")
	paths := strings.Split(PATH, ":")
	matches := []string{}
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.Mode().Perm()&0111 != 0 {
				if strings.HasPrefix(file.Name(), prefix) {
					matches = append(matches, file.Name())
				}
			}
		}
	}
	return matches
}

func Find(binary string) (bool, string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return false, "", fmt.Errorf("Not found")
	}
	return true, path, nil
}

func FindAll(prefix string) ([]string, error) {
	return scan(prefix), nil
}
