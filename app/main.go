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
			if bool := slices.Contains(commands, lines[1]); bool == true {
				fmt.Printf("%s is a shell builtin\n", lines[1])
			} else {
				fmt.Printf("%s: command not found\n", lines[1])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
