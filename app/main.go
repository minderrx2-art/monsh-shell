package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdout)
		// TODO: Uncomment the code below to pass the first stage
		fmt.Print("$ ")
		line, _ := reader.ReadString('\n')
		lines := strings.Split(strings.TrimSpace(line), " ")
		command, rest := lines[0], lines[1:]
		if command == "exit" {
			return
		}
		if command == "echo" {
			fmt.Printf("%s\n", strings.Join(rest, " "))
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
