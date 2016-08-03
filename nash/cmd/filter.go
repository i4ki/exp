package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	temp := make([]byte, 0, 1024)

	for {
		buf := make([]byte, 1024)
		_, err := os.Stdin.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Printf("err: %s\n", err.Error())
			os.Exit(1)
		}

		temp = append(temp, buf...)

	again:
		for i := 0; i < len(temp); i++ {
			if temp[i] == '\n' {
				line := temp[:i]
				if string(line) == "AB" || string(line) == "CD" {
					fmt.Println(string(line))
				}

				temp = temp[i+1:]

				break
			}

			if i == (len(temp) - 1) {
				temp = []byte{}
			}
		}

		if len(temp) > 0 {
			goto again
		}

		if err == io.EOF {
			break
		}
	}
}
