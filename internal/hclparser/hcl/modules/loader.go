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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/plancost/terraform-provider-plancost/internal/hclparser/config"
	"github.com/rs/zerolog"
)

var (
	// tfManifestPath is the name of the terraform module manifest file which stores the metadata of the modules
	tfManifestPath = ".terraform/modules/modules.json"
)

// ModuleLoader handles the loading of Terraform modules. It supports local, registry and other remote modules.
//
// The path should be the root directory of the Terraform project. We use a distinct module loader per Terraform project,
// because at the moment the cache is per project. The cache reads the manifest.json file from the path's
// .plancost/terraform_modules directory. We could implement a global cache in the future, but for now have decided
// to go with the same approach as Terraform.
type ModuleLoader struct {
	// cachePath is the path to the directory that Infracost will download modules to.
	// This is normally the top level directory of a multi-project environment, where the
	// Infracost config file resides or project auto-detection starts from.
	cachePath      string
	cache          *Cache
	hclParser      *SharedHCLParser
	sourceMap      config.TerraformSourceMap
	sourceMapRegex config.TerraformSourceMapRegex

	logger zerolog.Logger
}

type SourceMapResult struct {
	Source   string
	Version  string
	RawQuery string
}
type ModuleLoaderOptions struct {
	CachePath      string
	HCLParser      *SharedHCLParser
	SourceMap      config.TerraformSourceMap
	SourceMapRegex config.TerraformSourceMapRegex
	Logger         zerolog.Logger
}

// NewModuleLoader constructs a new module loader
func NewModuleLoader(opts ModuleLoaderOptions) *ModuleLoader {
	if err := opts.SourceMapRegex.Compile(); err != nil {
		opts.Logger.Error().Err(err).Msg("error compiling source map regex")
	}

	m := &ModuleLoader{
		cachePath:      opts.CachePath,
		cache:          NewCache(opts.Logger),
		hclParser:      opts.HCLParser,
		sourceMap:      opts.SourceMap,
		sourceMapRegex: opts.SourceMapRegex,
		logger:         opts.Logger,
	}

	return m
}

// tfManifestFilePath is the path to the Terraform module manifest file relative to the current working directory.
func (m *ModuleLoader) tfManifestFilePath(path string) string {
	return filepath.Join(path, tfManifestPath)
}

// Load loads the modules from the given path.
// For each module it checks if the module has already been downloaded, by checking if iut exists in the manifest
// If not then it downloads the module from the registry or from a remote source and updates the module manifest with the latest metadata.
func (m *ModuleLoader) Load(path string) (man *Manifest, err error) {
	defer func() {
		if man != nil {
			man.cachePath = m.cachePath
		}
	}()

	// Note-plancost: for now we only support reading from the Terraform module manifest file.
	tfManifestFilePath := m.tfManifestFilePath(path)
	_, err = os.Stat(tfManifestFilePath)
	if err == nil {
		manifest, err := readManifest(tfManifestFilePath)
		if err == nil {
			// let's make the module dirs relative to the path directory as later
			// we'll look up the modules based on the cache path at the Infracost root (where the infracost.yml
			// resides or where the --path autodetect started for multi-project)
			for i, module := range manifest.Modules {
				dir := path
				if m.cachePath != "" {
					dir, _ = filepath.Rel(m.cachePath, path)
				}

				manifest.Modules[i].Dir = filepath.Join(dir, module.Dir)
			}

			return manifest, nil
		}

		m.logger.Debug().Err(err).Msg("error reading terraform module manifest")
	}

	// return the root module if not found in the manifest
	return &Manifest{
		cachePath: "",
		Path:      path,
		Version:   "2",
		Modules:   nil,
	}, nil
}

func HasOpenTofuExtension(name string) bool {
	return filepath.Ext(name) == ".tofu" || strings.HasSuffix(name, ".tofu.json")
}

// IsLocalModule checks if the module is a local module by checking
// if the module source starts with any known local prefixes
func IsLocalModule(source string) bool {
	return strings.HasPrefix(source, "./") ||
		strings.HasPrefix(source, "../") ||
		strings.HasPrefix(source, ".\\") ||
		strings.HasPrefix(source, "..\\")
}

func splitModuleSubDir(moduleSource string) (string, string, error) {
	moduleAddr, submodulePath := getter.SourceDirSubdir(moduleSource)
	if strings.HasPrefix(submodulePath, "../") {
		return "", "", fmt.Errorf("invalid submodule path '%s'", submodulePath)
	}

	return moduleAddr, submodulePath, nil
}
