package builtin

import (
	"fmt"
	"os"
)

func Pwd() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", wd)
	return nil
}
