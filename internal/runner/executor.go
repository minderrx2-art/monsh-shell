package runner

import (
	"os"
	"os/exec"

	"github.com/minderrx2-art/monsh/internal/parser"
)

func Execute(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// &{[{echo [damn son] [{1 test.txt}]}]}
func ExecutePipeline(pipeline *parser.Pipeline) error {
	for _, command := range pipeline.Commands {
		cmd := exec.Command(command.Name, command.Args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		for _, redirect := range command.Redirects {
			switch redirect.Type {
			case parser.In: // <
				if err := redirectIn(cmd, redirect.Target); err != nil {
					return err
				}
			case parser.Out: // >
				if err := redirectOut(cmd, redirect.Target); err != nil {
					return err
				}
			case parser.OutErr:
				if err := redirectOutErr(cmd, redirect.Target); err != nil {
					return err
				}
			case parser.Append:
				if err := redirectAppend(cmd, redirect.Target); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func redirectIn(cmd *exec.Cmd, target string) error {
	file, err := os.Open(target)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdin = file
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func redirectOut(cmd *exec.Cmd, target string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdout = file
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func redirectOutErr(cmd *exec.Cmd, target string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stderr = file
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func redirectAppend(cmd *exec.Cmd, target string) error {
	file, err := os.OpenFile(
		target,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdout = file
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
