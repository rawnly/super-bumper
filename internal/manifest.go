package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type jsonManifest struct {
	Version string `json:"version"`
}

var tomlVersionRegex = regexp.MustCompile(`(?m)^version\s*=\s*"([^"]+)"`)

func DetectVersion(dir string) (string, error) {
	manifests := []struct {
		filename string
		parser   func([]byte) (string, error)
	}{
		{"package.json", parseJSON},
		{"composer.json", parseJSON},
		{"Cargo.toml", parseTOML},
	}

	for _, m := range manifests {
		path := filepath.Join(dir, m.filename)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		version, err := m.parser(data)
		if err != nil {
			continue
		}

		if version != "" {
			return version, nil
		}
	}

	return "", fmt.Errorf("no manifest file found with version field")
}

func parseJSON(data []byte) (string, error) {
	var manifest jsonManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	return manifest.Version, nil
}

func parseTOML(data []byte) (string, error) {
	matches := tomlVersionRegex.FindSubmatch(data)
	if matches == nil {
		return "", fmt.Errorf("no version found in TOML")
	}
	return string(matches[1]), nil
}
