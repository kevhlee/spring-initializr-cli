package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kevhlee/sprout/initializr"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newInputPrompt(element initializr.MetadataText, title string, value *string) *huh.Input {
	return huh.NewInput().
		Title(title).
		Value(value).
		CharLimit(80)
}

func newMultiSelectPrompt(element initializr.MetadataMultiSelect, title string, value *[]string, disableFunc func(initializr.MetadataValue) bool) *huh.MultiSelect[string] {
	options := []huh.Option[string]{}
	for _, hierarchy := range element.Values {
		for _, value := range hierarchy.Values {
			if disableFunc != nil && disableFunc(value) {
				continue
			}
			options = append(options, huh.NewOption(value.Name, value.Id))
		}
	}

	slices.SortFunc(options, func(a, b huh.Option[string]) int {
		return strings.Compare(a.Key, b.Key)
	})

	return huh.NewMultiSelect[string]().
		Title(title).
		Options(options...).
		Value(value)
}

func newSelectPrompt(element initializr.MetadataSelect, title string, value *string) *huh.Select[string] {
	options := make([]huh.Option[string], len(element.Values))
	for i, value := range element.Values {
		options[i] = huh.NewOption(value.Name, value.Id)
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value)
}

func run() error {
	metadata, err := initializr.FetchMetadata(initializr.DefaultUrl)
	if err != nil {
		return err
	}

	opts := initializr.NewDefaultOptions(metadata)

	form1 := huh.NewForm(
		huh.NewGroup(
			newSelectPrompt(metadata.Type, "Project:", &opts.Type),
			newSelectPrompt(metadata.Language, "Language:", &opts.Language),
			newSelectPrompt(metadata.BootVersion, "Spring Boot:", &opts.BootVersion),
		),
		huh.NewGroup(
			newInputPrompt(metadata.GroupId, "Group:", &opts.GroupId),
			newInputPrompt(metadata.ArtifactId, "Artifact:", &opts.ArtifactId),
			newInputPrompt(metadata.Name, "Name:", &opts.Name),
			newInputPrompt(metadata.Version, "Version:", &opts.Version),
			newInputPrompt(metadata.Description, "Description:", &opts.Description),
			newInputPrompt(metadata.PackageName, "Package name:", &opts.PackageName),
			newSelectPrompt(metadata.Packaging, "Packaging:", &opts.Packaging),
			newSelectPrompt(metadata.JavaVersion, "Java:", &opts.JavaVersion),
		),
	).WithTheme(huh.ThemeBase16())

	if err := form1.Run(); err != nil {
		return err
	}

	form2 := huh.NewForm(
		huh.NewGroup(
			newMultiSelectPrompt(
				metadata.Dependencies,
				"Dependencies:",
				&opts.Dependencies,
				func(value initializr.MetadataValue) bool {
					return WithinRange(opts.BootVersion, value.VersionRange)
				},
			),
		),
	).WithTheme(huh.ThemeBase16())

	if err := form2.Run(); err != nil {
		return err
	}

	err = initializr.GenerateProject(initializr.DefaultUrl, opts)
	if err == nil {
		fmt.Printf("Spring project '%s' created. 🌱\n", opts.Name)
	}
	return err
}

func runForm(form *huh.Form) error {
	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}
	}
	return err
}
