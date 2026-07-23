package builtin

import (
	"fmt"
)

func Complete(args []string) {
	fmt.Printf("%s\n", "complete: git: no completion specification")
}
