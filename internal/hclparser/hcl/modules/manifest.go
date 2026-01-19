// Copyright 2021 Infracost Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modules

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Manifest is a struct that represents the JSON found in the manifest.json file in the .infracost dir
// It is used for caching the modules that have already been downloaded.
// It uses the same format as the .terraform/modules/modules.json file
type Manifest struct {
	cachePath string

	Path    string            `json:"Path"`
	Version string            `json:"Version"`
	Modules []*ManifestModule `json:"Modules"`
}

func (m Manifest) Get(key string) ManifestModule {
	for _, module := range m.Modules {
		if module.Key == key {
			loc := filepath.Clean(filepath.Join(m.cachePath, module.Dir))
			return ManifestModule{
				Key:         module.Key,
				Source:      module.Source,
				Version:     module.Version,
				Dir:         loc,
				DownloadURL: module.DownloadURL,
			}
		}
	}

	return ManifestModule{}
}

// ManifestModule represents a single module in the manifest.json file
type ManifestModule struct {
	Key            string `json:"Key"`
	Source         string `json:"Source"`
	Version        string `json:"Version,omitempty"`
	Dir            string `json:"Dir"`
	DownloadURL    string
	IsSourceMapped bool `json:"-"`
}

func (m ManifestModule) URL() string {
	if IsLocalModule(m.Source) {
		return ""
	}

	remoteSource := m.Source

	if m.DownloadURL != "" {
		remoteSource = m.DownloadURL
	}

	remoteSource, _, _ = splitModuleSubDir(remoteSource)
	remoteSource = strings.TrimPrefix(remoteSource, "git::")
	remoteSource = strings.TrimPrefix(remoteSource, "gcs::")
	remoteSource = strings.TrimPrefix(remoteSource, "s3::")

	u, err := url.Parse(remoteSource)
	if err == nil {
		u.RawQuery = ""
		return u.String()
	}

	return remoteSource
}

// readManifest reads the manifest file from the given path
func readManifest(path string) (*Manifest, error) {
	var manifest Manifest

	data, err := os.ReadFile(path)
	if err != nil {
		return &manifest, fmt.Errorf("failed to read module manifest: %w", err)
	}

	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return &manifest, fmt.Errorf("failed to unmarshal module manifest: %w", err)
	}

	return &manifest, err
}
