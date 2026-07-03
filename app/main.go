package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		command, second_command, rest := lines[0], "", []string{}

		if len(lines) > 1 {
			second_command, rest = lines[1], lines[1:]
		}

		if command == "exit" {
			return

		} else if command == "echo" {
			fmt.Printf("%s\n", strings.Join(rest, " "))

		} else if command == "type" {
			if contains := slices.Contains(commands, second_command); contains == true {
				fmt.Printf("%s is a shell builtin\n", second_command)

			} else if exists, path, err := findPath(second_command); exists == true && err == nil {
				fmt.Printf("%s is %s\n", second_command, path)

			} else {
				fmt.Printf("%s: not found\n", second_command)

			}
		} else if exists, path, err := findPath(command); exists == true && err == nil {
			runExecutable(path, rest...)

		} else if stats, err := os.Stat(command); err == nil {
			mode := stats.Mode()
			if mode&0111 != 0 {
				runExecutable(command)
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func findPath(binary string) (bool, string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return false, "", fmt.Errorf("Not found")
	}
	return true, path, nil
}

func runExecutable(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
