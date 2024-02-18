package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	qualifiers          = []string{"M", "RC", "BUILD-SNAPSHOT", "RELEASE"}
	patternVersion      = regexp.MustCompile(`^(\d+)\.(\d+|x)\.(\d+|x)(?:([.|-])([^0-9]+)(\d+)?)?$`)
	patternVersionRange = regexp.MustCompile(`(\(|\[)(.*),(.*)(\)|\])`)
)

func CompareQualifier(a, b string) int {
	indexA := len(qualifiers) - 1
	indexB := len(qualifiers) - 1

	for i := 0; i < len(qualifiers); i++ {
		if a == qualifiers[i] {
			indexA = i
		}
		if b == qualifiers[i] {
			indexB = i
		}
	}

	return indexA - indexB
}

func CompareVersion(a, b string) int {
	versionA := strings.Split(a, ".")
	versionB := strings.Split(b, ".")

	for i := 0; i < 3; i++ {
		va, _ := strconv.ParseInt(versionA[i], 10, 32)
		vb, _ := strconv.ParseInt(versionB[i], 10, 32)

		diff := va - vb
		if diff != 0 {
			return int(diff)
		}
	}

	return CompareQualifier(ParseQualifier(a), ParseQualifier(b))
}

func ParseQualifier(version string) string {
	match := patternVersion.FindStringSubmatch(version)
	if len(match) == 0 {
		return "RELEASE"
	}

	qualifier := match[5]
	if len(match) == 0 {
		return "RELEASE"
	}

	for _, q := range qualifiers {
		if q == qualifier {
			return qualifier
		}
	}
	return "RELEASE"
}

func WithinVersionRange(version, versionRange string) bool {
	if len(versionRange) == 0 {
		return true
	}

	match := patternVersionRange.FindStringSubmatch(versionRange)
	if len(match) == 0 {
		return false
	}

	var (
		lowerInclusive  = match[1] == "["
		lowerVersion    = match[2]
		higherVersion   = match[3]
		higherInclusive = match[4] == "]"
	)

	if lowerInclusive && higherInclusive {
		return CompareVersion(lowerVersion, version) <= 0 && CompareVersion(higherVersion, version) >= 0
	} else if lowerInclusive {
		return CompareVersion(lowerVersion, version) <= 0 && CompareVersion(higherVersion, version) > 0
	} else if higherInclusive {
		return CompareVersion(lowerVersion, version) < 0 && CompareVersion(higherVersion, version) >= 0
	} else {
		return CompareVersion(lowerVersion, version) < 0 && CompareVersion(higherVersion, version) > 0
	}
}
