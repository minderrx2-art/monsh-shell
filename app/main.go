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
		line = strings.TrimSpace(line)

		command, rest := strings.Split(line, " ")[0], strings.Split(line, " ")[1:]
		if command == "exit" {
			return
		}
		if command == "echo" {
			fmt.Printf("%s\n", rest)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
