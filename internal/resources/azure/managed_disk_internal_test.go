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

package azure

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestMapDiskName(t *testing.T) {
	tests := []struct {
		diskType      string
		requestedSize int
		expected      string
	}{
		{"Standard", 32, "S4"},
		{"Standard", 4000, "S50"},
		{"Standard", 8192, "S60"},
		{"Standard", 400000, ""},
		{"StandardSSD", 32, "E4"},
		{"StandardSSD", 4000, "E50"},
		{"StandardSSD", 8192, "E60"},
		{"StandardSSD", 400000, ""},
		{"Premium", 32, "P4"},
		{"Premium", 4000, "P50"},
		{"Premium", 8192, "P60"},
		{"Premium", 400000, ""},
	}

	for _, test := range tests {
		actual := mapDiskName(test.diskType, test.requestedSize)
		assert.Equal(t, test.expected, actual)
	}
}

func TestMapUltraDiskSize(t *testing.T) {
	tests := []struct {
		requestedSize int
		expected      int
	}{
		{0, 4},
		{3, 4},
		{4, 4},
		{5, 8},
		{1023, 1024},
		{1024, 1024},
		{1025, 2048},
		{65536, 65536},
		{65537, 65536},
	}

	for _, test := range tests {
		actual := mapUltraDiskSize(test.requestedSize)
		assert.Equal(t, test.expected, actual)
	}
}

func TestMapStorageReplicationType(t *testing.T) {
	tests := []struct {
		storageType string
		expected    bool
	}{
		{"LRS", true},
		{"ZRS", true},
		{"ABC", false},
	}

	for _, test := range tests {
		actual := mapStorageReplicationType(test.storageType)
		assert.Equal(t, test.expected, actual)
	}
}
