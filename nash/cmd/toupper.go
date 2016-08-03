package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	buf := make([]byte, 1024)

	for {
		n, err := os.Stdin.Read(buf)

		if err == io.EOF || n == 0 {
			break
		}

		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			break
		}
	}

	fmt.Printf("%s\n", strings.ToUpper(string(buf)))
}
