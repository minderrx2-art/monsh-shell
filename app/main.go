package main

import (
	"fmt"
)

func main() {
	for {
		var command string
		// TODO: Uncomment the code below to pass the first stage
		fmt.Print("$ ")
		fmt.Scan(&command)
		if command == "exit" {
			return
		}
		fmt.Printf("%s: command not found\n", command)
	}
}
