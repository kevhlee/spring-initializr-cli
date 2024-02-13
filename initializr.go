package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// FetchMetadata fetches JSON metadata containing settings that can be used to
// generate a project from a Spring Initializr API instance.
//
// For reference: https://docs.spring.io/initializr/docs/current/reference/html/#api-guide
func FetchMetadata(urlpath string) (Metadata, error) {
	var metadata Metadata

	req, err := http.NewRequest(http.MethodGet, urlpath, nil)
	if err != nil {
		return metadata, err
	}
	req.Header.Add("Accept", "application/vnd.initializr.v2.2+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return metadata, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return metadata, err
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return metadata, err
	}
	return metadata, nil
}

// Metadata is JSON metadata containing settings that can be used to create a
// new project from a Spring Initializr API instance.
type Metadata struct {
	ArtifactId   MetadataText        `json:"artifactId"`
	BootVersion  MetadataSelect      `json:"bootVersion"`
	Dependencies MetadataMultiSelect `json:"dependencies"`
	Description  MetadataText        `json:"description"`
	GroupId      MetadataText        `json:"groupId"`
	JavaVersion  MetadataSelect      `json:"javaVersion"`
	Language     MetadataSelect      `json:"language"`
	Name         MetadataText        `json:"name"`
	PackageName  MetadataText        `json:"packageName"`
	Packaging    MetadataSelect      `json:"packaging"`
	Type         MetadataSelect      `json:"type"`
	Version      MetadataText        `json:"version"`
}

type MetadataType string

type MetadataHierarchy struct {
	Name   string          `json:"name"`
	Values []MetadataValue `json:"values"`
}

type MetadataMultiSelect struct {
	Type   MetadataType        `json:"type"`
	Values []MetadataHierarchy `json:"values"`
}

type MetadataSelect struct {
	Type    MetadataType    `json:"type"`
	Default string          `json:"default"`
	Values  []MetadataValue `json:"values"`
}

type MetadataText struct {
	Type    MetadataType `json:"type"`
	Default string       `json:"default"`
}

type MetadataValue struct {
	Description  string `json:"description"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	VersionRange string `json:"versionRange"`
}

func (v MetadataValue) HasDescription() bool {
	return len(v.Description) != 0
}

func (v MetadataValue) HasVersionRange() bool {
	return len(v.VersionRange) != 0
}

// GenerateProject sends a request to a Spring Initializr instance to create
// a project based on specified options.
func GenerateProject(urlpath string, opts Options) error {
	u, err := url.Parse(urlpath)
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
	req.Header.Add("Accept", "application/vnd.initializr.v2.2+json")

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
