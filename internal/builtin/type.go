package builtin

import (
	"fmt"
	"slices"

	"github.com/minderrx2-art/monsh/internal/path"
)

var builtin = []string{
	"cd",
	"exit",
	"echo",
	"type",
	"pwd",
}

func Type(second_command string) {
	if isBuiltin := slices.Contains(builtin, second_command); isBuiltin == true {
		fmt.Printf("%s is a shell builtin\n", second_command)

	} else if exists, path, err := path.Find(second_command); exists == true && err == nil {
		fmt.Printf("%s is %s\n", second_command, path)

	} else {
		fmt.Printf("%s: not found\n", second_command)

	}
}
