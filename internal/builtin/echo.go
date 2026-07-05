package builtin

import (
	"fmt"
	"strings"
)

func Echo(input []string) {
	fmt.Printf("%s\n", strings.Join(input, " "))
}
