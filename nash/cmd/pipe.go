// Example program to execute a pipeline of programs writing
// the result to a file.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type (
	// Cmd is the command specification
	Cmd struct {
		*exec.Cmd

		fdIn  map[uint]io.Reader
		fdOut map[uint]io.Writer
	}
)

func NewCmd(path string) Cmd {
	cmd := Cmd{}
	cmd.Cmd = &exec.Cmd{}
	cmd.Cmd.Path = path
	cmd.fdIn = make(map[uint]io.Reader)
	cmd.fdOut = make(map[uint]io.Writer)

	return cmd
}

func (c *Cmd) SetInputfd(n uint, in io.Reader) {
	c.fdIn[n] = in
}

func (c *Cmd) SetOutputfd(n uint, out io.Writer) {
	c.fdOut[n] = out
}

func (c *Cmd) setStdin(value io.Reader) {
	c.Cmd.Stdin = value
}

func (c *Cmd) setStdout(value io.Writer) {
	c.Cmd.Stdout = value
}

func (c *Cmd) setStderr(value io.Writer) {
	c.Cmd.Stderr = value
}

func (c *Cmd) addExtraFile(value *os.File) {
	if c.Cmd.ExtraFiles == nil {
		c.Cmd.ExtraFiles = make([]*os.File, 0, 8)
	}

	c.Cmd.ExtraFiles = append(c.Cmd.ExtraFiles, value)
}

func (c *Cmd) applyFd() error {
	for fd, value := range c.fdIn {
		switch fd {
		case 0:
			c.setStdin(value)
		default:
			file, ok := value.(*os.File)

			if !ok {
				return fmt.Errorf("ExtraFiles requires a file object.")
			}

			c.addExtraFile(file)
		}
	}

	for fd, value := range c.fdOut {
		switch fd {
		case 1:
			c.setStdout(value)
		case 2:
			c.setStderr(value)
		default:
			file, ok := value.(*os.File)

			if !ok {
				return fmt.Errorf("ExtraFiles requires a file object.")
			}

			c.addExtraFile(file)
		}
	}

	return nil
}

func (c *Cmd) Start() error {
	err := c.applyFd()

	if err != nil {
		return err
	}

	return c.Cmd.Start()
}

func (c *Cmd) Wait() error {
	return c.Cmd.Wait()
}

func buildPipe(cmds []*Cmd) error {
	var err error

	if len(cmds) < 2 {
		return errors.New("at least 2 programs")
	}

	last := len(cmds) - 1

	cmds[0].SetInputfd(0, os.Stdin)

	for i := 0; i < last; i++ {
		cmd := cmds[i]

		cmd.SetOutputfd(2, os.Stderr)

		stdin, err := cmd.StdoutPipe()

		if err != nil {
			return err
		}

		cmds[i+1].SetInputfd(0, stdin)
	}

	file, err := os.OpenFile("./out", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	cmds[last].SetOutputfd(2, os.Stderr)
	cmds[last].SetOutputfd(1, file)

	return nil
}

func main() {
	var err error

	cmd1 := NewCmd("./gen")
	cmd2 := NewCmd("./toupper")
	cmd3 := NewCmd("./filter")

	cmds := []*Cmd{
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
