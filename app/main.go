package main

import (
	"fmt"
)

func main() {
	var command string
	// TODO: Uncomment the code below to pass the first stage
	fmt.Print("$ ")
	fmt.Scan(&command)
	fmt.Printf("%s: command not found\n", command)
}
