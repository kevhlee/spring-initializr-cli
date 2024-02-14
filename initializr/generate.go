package initializr

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// GenerateProject sends a request to a Spring Initializr instance to create
// a project based on specified options.
func GenerateProject(rawURL string, opts Options) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	u = u.JoinPath("starter.zip")

	values := u.Query()
	values.Set("artifactId", opts.ArtifactId)
	values.Set("bootVersion", opts.BootVersion)
	values.Set("description", opts.Description)
	values.Add("dependencies", strings.Join(opts.Dependencies, ","))
	values.Set("groupId", opts.GroupId)
	values.Set("javaVersion", opts.JavaVersion)
	values.Set("language", opts.Language)
	values.Set("packageName", opts.PackageName)
	values.Set("packaging", opts.Packaging)
	values.Set("name", opts.Name)
	values.Set("type", opts.Type)
	values.Set("version", opts.Version)

	u.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", DefaultAcceptHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	basePath := opts.Name

	for _, f := range zr.File {
		filename := path.Join(basePath, f.Name)
		fileinfo := f.FileInfo()

		if fileinfo.IsDir() {
			if err := os.MkdirAll(filename, fileinfo.Mode()); err != nil {
				return err
			}
		} else {
			reader, err := f.Open()
			if err != nil {
				return err
			}

			data, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filename, data, fileinfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Options represents options for Spring Initializr project creation.
type Options struct {
	ArtifactId   string
	BootVersion  string
	Dependencies []string
	Description  string
	GroupId      string
	JavaVersion  string
	Language     string
	Name         string
	PackageName  string
	Packaging    string
	Type         string
	Version      string
}

// NewDefaultOptions creates options based on Spring Initializr metadata.
func NewDefaultOptions(metadata Metadata) Options {
	return Options{
		ArtifactId:  metadata.ArtifactId.Default,
		BootVersion: metadata.BootVersion.Default,
		Description: metadata.Description.Default,
		GroupId:     metadata.GroupId.Default,
		JavaVersion: metadata.JavaVersion.Default,
		Language:    metadata.Language.Default,
		PackageName: metadata.PackageName.Default,
		Packaging:   metadata.Packaging.Default,
		Name:        metadata.Name.Default,
		Type:        metadata.Type.Default,
		Version:     metadata.Version.Default,
	}
}
