package builtin

import (
	"fmt"
)

func Complete(args []string) {
	switch args[0] {
	case "-p":
		fmt.Printf("complete: %s: no completion specification\n", args[1])
	default:
		fmt.Printf("complete: %s: no completion specification\n", args[1])
	}
}
