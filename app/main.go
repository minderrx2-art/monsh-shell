package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	commands := []string{"exit", "echo", "type"}
	for {
		reader := bufio.NewReader(os.Stdout)
		fmt.Print("$ ")
		line, _ := reader.ReadString('\n')
		lines := strings.Split(strings.TrimSpace(line), " ")
		command, rest := lines[0], lines[1:]

		if command == "exit" {
			return
		} else if command == "echo" {
			fmt.Printf("%s\n", strings.Join(rest, " "))
		} else if command == "type" {
			if contains := slices.Contains(commands, lines[1]); contains == true {
				fmt.Printf("%s is a shell builtin\n", lines[1])
			} else if exists, path, err := readBinPath(lines[1]); exists == true && err == nil {
				fmt.Printf("%s is %s\n", lines[1], path)
			} else {
				fmt.Printf("%s: not found\n", lines[1])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err == os.ErrNotExist {
		return false
	}
	return true
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	if mode.IsRegular() {
		return false
	}
	isUserExecutable := mode&0100 != 0
	return isUserExecutable
}

func readBinPath(binary string) (bool, string, error) {
	paths := strings.Split(os.Getenv("PATH"), ":")
	filtered_paths := slices.DeleteFunc(paths, func(path string) bool {
		if strings.Contains(path, binary) {
			return false
		}
		return true
	})
	for _, f_path := range filtered_paths {
		if fileExists(f_path) && isExecutable(f_path) {
			return true, f_path, nil
		}
	}
	return false, "", fmt.Errorf("Binary has no PATH")
}
