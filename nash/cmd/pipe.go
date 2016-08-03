// Example program to execute a pipeline of programs writing
// the result to a file.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func buildPipe(cmds []*exec.Cmd) error {
	var err error

	if len(cmds) < 2 {
		return errors.New("at least 2 programs")
	}

	last := len(cmds) - 1

	cmds[0].Stdin = os.Stdin

	for i := 0; i < last; i++ {
		cmd := cmds[i]

		cmd.Stderr = os.Stderr
		cmds[i+1].Stdin, err = cmd.StdoutPipe()

		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile("./out", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	cmds[last].Stderr = os.Stderr
	cmds[last].Stdout = file

	return nil
}

func main() {
	var err error

	cmd1 := exec.Cmd{
		Path: "./gen",
	}

	cmd2 := exec.Cmd{
		Path: "./toupper",
	}

	cmd3 := exec.Cmd{
		Path: "./filter",
	}

	cmds := []*exec.Cmd{
		&cmd1,
		&cmd2,
		&cmd3,
	}

	err = buildPipe(cmds)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	for _, cmd := range cmds {
		fmt.Printf("starting: %s\n", cmd.Path)
		if err = cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "start: %s\n", err.Error())
			os.Exit(1)
		}
	}

	for _, cmd := range cmds {
		fmt.Printf("Waiting: %s\n", cmd.Path)
		if err = cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "wait: %s\n", err.Error())
			os.Exit(1)
		}
	}
}
