package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	if err := startCLI("https://start.spring.io"); err != nil {
		if err != huh.ErrUserAborted {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func startCLI(urlPath string) error {
	metadata, err := FetchMetadata(urlPath)
	if err != nil {
		return err
	}

	opts, err := StartPrompts(metadata)
	if err != nil {
		return err
	}

	if err := GenerateProject(urlPath, opts); err != nil {
		return err
	}

	fmt.Printf("Spring project '%s' created. 🌱\n", opts.Name)
	return nil
}
