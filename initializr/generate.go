package initializr

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// GenerateProject sends a request to a Spring Initializr instance to create
// a project based on specified options.
func GenerateProject(options Options) error {
	u, err := url.Parse(DefaultUrl)
	if err != nil {
		return err
	}
	u = u.JoinPath("starter.zip")

	values := u.Query()
	values.Set("artifactId", options.ArtifactId)
	values.Set("bootVersion", options.BootVersion)
	values.Set("description", options.Description)
	values.Add("dependencies", strings.Join(options.Dependencies, ","))
	values.Set("groupId", options.GroupId)
	values.Set("javaVersion", options.JavaVersion)
	values.Set("language", options.Language)
	values.Set("packageName", options.PackageName)
	values.Set("packaging", options.Packaging)
	values.Set("name", options.Name)
	values.Set("type", options.Type)
	values.Set("version", options.Version)

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

	basePath := options.Name
	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		return fmt.Errorf("Directory '%s' already exists", basePath)
	}

	if err := os.MkdirAll(basePath, 0777); err != nil {
		return err
	}

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
	ArtifactId   string   `json:"artifactId"`
	BootVersion  string   `json:"bootVersion"`
	Dependencies []string `json:"dependencies"`
	Description  string   `json:"description"`
	GroupId      string   `json:"groupId"`
	JavaVersion  string   `json:"javaVersion"`
	Language     string   `json:"language"`
	Name         string   `json:"name"`
	PackageName  string   `json:"packageName"`
	Packaging    string   `json:"packaging"`
	Type         string   `json:"type"`
	Version      string   `json:"version"`
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

// ParseOptions parses options from a JSON file.
func ParseOptions(filename string) (options Options, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return options, err
	}

	if err := json.Unmarshal(data, &options); err != nil {
		return options, err
	}
	return options, nil
}
