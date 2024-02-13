package main

import (
	"fmt"
	"os"
)

func main() {
	if err := StartCLI(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
