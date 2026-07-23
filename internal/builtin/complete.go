package builtin

import (
	"fmt"
)

func Complete(args []string) {
	println([]string(args))
	switch args[0] {
	case "-p":
		fmt.Printf("complete: %s: no completion specification", args[1])
	default:
		fmt.Printf("complete: %s: no completion specification", args[1])
	}
}
