package path

import (
	"fmt"
	"os"
	"os/exec"
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

func scanWd() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return []string{}, err
	}
	matches, err := os.ReadDir(currentDir)
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
	var splitPrefix string = strings.Split(prefix, " ")[1]
	files, err := scanWd()
	fmt.Println("files", files)
	if err != nil {
		return []string{}
	}
	matches := slices.DeleteFunc(files, func(file string) bool {
		return !strings.HasPrefix(file, splitPrefix)
	})
	fmt.Println("matches", matches)
	return matches
}
