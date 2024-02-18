package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kevhlee/sprout/initializr"
	"github.com/spf13/cobra"
)

func main() {
	var (
		filename string
	)

	cmd := &cobra.Command{
		Use:           "sprout",
		Short:         "Sprout CLI 🌱",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(filename)
		},
	}

	cmd.Flags().StringVarP(&filename, "file", "f", "", "path to file containing Spring Initializr options")

	cobra.CheckErr(cmd.Execute())
}

func run(filename string) (err error) {
	var options initializr.Options

	if len(filename) == 0 {
		options, err = runPrompts()
	} else {
		options, err = initializr.ParseOptions(filename)
	}

	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	if err := initializr.GenerateProject(options); err != nil {
		return err
	}

	fmt.Printf("Spring project '%s' generated. 🌱\n", options.Name)
	return nil
}

func runPrompts() (options initializr.Options, err error) {
	metadata, err := initializr.FetchMetadata()
	if err != nil {
		return options, err
	}

	options = initializr.NewDefaultOptions(metadata)

	err = runForm(
		huh.NewGroup(
			newSelectPrompt("Project:", &options.Type, metadata.Type),
			newSelectPrompt("Language:", &options.Language, metadata.Language),
			newSelectPrompt("Spring Boot:", &options.BootVersion, metadata.BootVersion),
		),
		huh.NewGroup(
			newInputPrompt("Group:", &options.GroupId),
			newInputPrompt("Artifact:", &options.ArtifactId),
			newInputPrompt("Name:", &options.Name),
			newInputPrompt("Description:", &options.Description),
			newInputPrompt("Package name:", &options.PackageName),
			newSelectPrompt("Packaging:", &options.Packaging, metadata.Packaging),
			newSelectPrompt("Java:", &options.JavaVersion, metadata.JavaVersion),
		),
	)

	if err != nil {
		return options, err
	}

	dependencies := []huh.Option[string]{}
	for _, hierarchy := range metadata.Dependencies.Values {
		for _, value := range hierarchy.Values {
			if WithinVersionRange(options.BootVersion, value.VersionRange) {
				dependencies = append(dependencies, huh.NewOption(value.Name, value.Id))
			}
		}
	}

	slices.SortFunc(dependencies, func(a, b huh.Option[string]) int {
		return strings.Compare(a.Key, b.Key)
	})

	err = runForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Dependencies:").
				Options(dependencies...).
				Value(&options.Dependencies),
		),
	)

	if err != nil {
		return options, err
	}

	var confirm bool = true

	err = runForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate project?").
				Value(&confirm).
				Description(getOptionsDescription(options)).
				Affirmative("Yes").
				Negative("No"),
		),
	)

	if err != nil {
		return options, err
	}

	if !confirm {
		err = huh.ErrUserAborted
	}
	return options, err
}

func runForm(groups ...*huh.Group) error {
	err := huh.NewForm(groups...).WithTheme(huh.ThemeBase16()).Run()
	if err != nil {
		return err
	}
	return nil
}

func newInputPrompt(title string, value *string) *huh.Input {
	return huh.NewInput().
		Title(title).
		Value(value).
		CharLimit(100)
}

func newSelectPrompt(title string, value *string, element initializr.MetadataSelect) *huh.Select[string] {
	options := make([]huh.Option[string], len(element.Values))
	for i, value := range element.Values {
		options[i] = huh.NewOption(value.Name, value.Id)
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value)
}

func getOptionsDescription(options initializr.Options) string {
	description := strings.Builder{}

	description.WriteString(fmt.Sprintf("%12s: %s\n", "Project", options.Type))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Language", options.Language))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Spring Boot", options.BootVersion))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Group", options.GroupId))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Artifact", options.ArtifactId))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Name", options.Name))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Description", options.Description))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Package name", options.PackageName))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Packaging", options.Packaging))
	description.WriteString(fmt.Sprintf("%12s: %s\n", "Java", options.JavaVersion))

	description.WriteString(fmt.Sprintf("%12s:\n", "Dependencies"))
	for _, dependency := range options.Dependencies {
		description.WriteString(fmt.Sprintf("%4s- %s\n", " ", dependency))
	}

	return description.String()
}
