package path

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
)

func scanPath() []string {
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
				matches = append(matches, file.Name())
			}
		}
	}
	return matches
}

func scanPwd() ([]string, error) {
	matches, err := os.ReadDir(path.Join(os.Getenv("PWD")))
	if err != nil {
		return []string{}, err
	}
	list := []string{}
	for _, match := range matches {
		list = append(list, match.Name())
	}
	return list, nil
}

func FindExecutable(binary string) (string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return "", fmt.Errorf("Executable not found")
	}
	return path, nil
}

func ListExecutables(prefix string) []string {
	executables := scanPath()
	matches := slices.DeleteFunc(executables, func(executable string) bool {
		return !strings.HasPrefix(executable, prefix)
	})

	if len(matches) == 0 {
		return []string{}
	}
	return matches
}

func ListFiles(prefix string) []string {
	files, err := scanPwd()
	if err != nil {
		return []string{}
	}
	matches := slices.DeleteFunc(files, func(file string) bool {
		return !strings.HasPrefix(file, prefix)
	})
	return matches
}
