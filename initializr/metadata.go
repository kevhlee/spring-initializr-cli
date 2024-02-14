package initializr

import (
	"encoding/json"
	"io"
	"net/http"
)

// FetchMetadata fetches JSON metadata containing settings that can be used to
// generate a project from a Spring Initializr API instance.
//
// For reference: https://docs.spring.io/initializr/docs/current/reference/html/#api-guide
func FetchMetadata(rawURL string) (Metadata, error) {
	var metadata Metadata

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return metadata, err
	}
	req.Header.Add("Accept", DefaultAcceptHeader)

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
