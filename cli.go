package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
)

func StartCLI() error {
	// TODO: Add command-line parsing
	initializrUrlPath := "https://start.spring.io"

	metadata, err := FetchMetadata(initializrUrlPath)
	if err != nil {
		return err
	}

	opts := NewDefaultOptions(metadata)

	form := huh.NewForm(
		huh.NewGroup(
			NewSelectPrompt(metadata.Type, "Project:", &opts.Type, nil),
			NewSelectPrompt(metadata.Language, "Language:", &opts.Language, nil),
			NewSelectPrompt(metadata.BootVersion, "Spring Boot:", &opts.BootVersion, nil),
		),
		huh.NewGroup(
			NewInputPrompt(metadata.GroupId, "Group:", &opts.GroupId),
			NewInputPrompt(metadata.ArtifactId, "Artifact:", &opts.ArtifactId),
			NewInputPrompt(metadata.Name, "Name:", &opts.Name),
			NewInputPrompt(metadata.Version, "Version:", &opts.Version),
			NewInputPrompt(metadata.Description, "Description:", &opts.Description),
			NewInputPrompt(metadata.PackageName, "Package name:", &opts.PackageName),
			NewSelectPrompt(metadata.Packaging, "Packaging:", &opts.Packaging, nil),
			NewSelectPrompt(metadata.JavaVersion, "Java:", &opts.JavaVersion, nil),
		),
		huh.NewGroup(
			NewMultiSelectPrompt(
				metadata.Dependencies,
				"Dependencies:",
				&opts.Dependencies,
				func(value MetadataValue) bool {
					return WithinRange(opts.BootVersion, value.VersionRange)
				},
			),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// TODO: Add progress bar
	if err := GenerateProject(initializrUrlPath, opts); err != nil {
		return err
	}

	fmt.Printf("Spring project '%s' created. 🌱\n", opts.Name)
	return nil
}

func NewInputPrompt(element MetadataText, title string, value *string) *huh.Input {
	return huh.NewInput().
		Title(title).
		Value(value).
		CharLimit(80)
}

func NewMultiSelectPrompt(element MetadataMultiSelect, title string, value *[]string, disableFunc func(MetadataValue) bool) *huh.MultiSelect[string] {
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

func NewSelectPrompt(element MetadataSelect, title string, value *string, disableFunc func(MetadataValue) bool) *huh.Select[string] {
	options := make([]huh.Option[string], len(element.Values))
	for i, value := range element.Values {
		if disableFunc != nil && disableFunc(value) {
			continue
		}
		options[i] = huh.NewOption(value.Name, value.Id)
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value)
}
