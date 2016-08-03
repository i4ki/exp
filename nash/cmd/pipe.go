// Example program to execute a pipeline of programs writing
// the result to a file.
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var err, err1, err2 error

	cmd1 := exec.Cmd{
		Path: "./gen",
	}

	cmd2 := exec.Cmd{
		Path: "./toupper",
	}

	cmd3 := exec.Cmd{
		Path: "./filter",
	}

	cmd1.Stdin = os.Stdin
	cmd1.Stderr = os.Stderr

	cmd2.Stdin, err1 = cmd1.StdoutPipe()
	cmd2.Stderr = os.Stderr

	file, err := os.OpenFile("./out", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	cmd3.Stdin, err2 = cmd2.StdoutPipe()
	cmd3.Stderr = os.Stderr
	cmd3.Stdout = file

	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Err: %s, %s\n", err1, err2)
		os.Exit(1)
	}

	if err = cmd1.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if err = cmd2.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if err = cmd3.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if err = cmd1.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if err = cmd2.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if err = cmd3.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

}
