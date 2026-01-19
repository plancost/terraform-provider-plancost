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
	"errors"
	"fmt"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog"
)

// Cache is a cache of modules that can be used to lookup modules to check if they've already been loaded.
//
// This only works with modules that have the same identifier. It doesn't cache modules that are used
// multiple times with different identifiers. That is done separately by the PackageFetcher and only
// caches per-run of Infracost, so if you add the same module to your Terraform code it will redownload that module.
// We could optimize it by moving the package fetching cache logic into here, but it would be inconsistent
// with how terraform init works.
type Cache struct {
	keyMap sync.Map
	logger zerolog.Logger
}

// NewCache creates a new cache from a module manifest
func NewCache(logger zerolog.Logger) *Cache {
	return &Cache{
		logger: logger,
	}
}

// lookupModule looks up a module in the cache by its key and checks that the
// source and version are compatible with the module in the cache.
func (c *Cache) lookupModule(key string, moduleCall *tfconfig.ModuleCall) (*ManifestModule, error) {
	value, ok := c.keyMap.Load(key)
	if !ok {
		return nil, errors.New("not in cache")
	}
	manifestModule, _ := value.(*ManifestModule)

	if manifestModule.Source == moduleCall.Source {
		return checkVersion(moduleCall, manifestModule)
	}

	return nil, errors.New("source has changed")
}

func checkVersion(moduleCall *tfconfig.ModuleCall, manifestModule *ManifestModule) (*ManifestModule, error) {
	if moduleCall.Version != "" && manifestModule.Version != "" {
		constraints, err := goversion.NewConstraint(moduleCall.Version)
		if err != nil {
			return nil, fmt.Errorf("invalid version constraint: %w", err)
		}

		version, err := goversion.NewVersion(manifestModule.Version)
		if err != nil {
			return nil, fmt.Errorf("invalid version: %w", err)
		}

		if !constraints.Check(version) {
			return nil, errors.New("version constraint doesn't match")
		}
	}

	return manifestModule, nil
}
