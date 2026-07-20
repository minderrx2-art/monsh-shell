package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/minderrx2-art/monsh/internal/parser"
	"github.com/minderrx2-art/monsh/internal/path"
)

func Execute(execPath string, args ...string) error {
	cmd := exec.Command(execPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil
		}
		return err
	}
	return nil
}

func ExecutePipeline(pipeline *parser.Pipeline) error {
	for _, command := range pipeline.Commands {
		if _, err := path.FindExecutable(command.Name); err != nil {
			return fmt.Errorf("%s: command not found", command.Name)
		}

		cmd := exec.Command(command.Name, command.Args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		var files []*os.File
		var err error
		for _, redirect := range command.Redirects {
			var file *os.File
			switch redirect.Type {
			case parser.In:
				file, err = openInput(redirect.Target)
				if err == nil {
					cmd.Stdin = file
				}
			case parser.Out:
				file, err = createOutput(redirect.Target)
				if err == nil {
					cmd.Stdout = file
				}
			case parser.OutErr:
				file, err = createOutput(redirect.Target)
				if err == nil {
					cmd.Stderr = file
				}
			case parser.Append:
				file, err = openAppend(redirect.Target)
				if err == nil {
					cmd.Stdout = file
				}
			case parser.AppendErr:
				file, err = openAppend(redirect.Target)
				if err == nil {
					cmd.Stderr = file
				}
			}
			if err != nil {
				for _, f := range files {
					f.Close()
				}
				return err
			}
			if file != nil {
				files = append(files, file)
			}
		}

		if err := cmd.Run(); err != nil {
			for _, f := range files {
				f.Close()
			}
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				continue
			}
			return err
		}
		for _, f := range files {
			f.Close()
		}
	}
	return nil
}

func openInput(target string) (*os.File, error) {
	file, err := os.Open(target)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createOutput(target string) (*os.File, error) {
	return os.Create(target)
}

func openAppend(target string) (*os.File, error) {
	return os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}
