package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			return
		default:
			var command string
			// TODO: Uncomment the code below to pass the first stage
			fmt.Print("$ ")
			fmt.Scan(&command)
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
