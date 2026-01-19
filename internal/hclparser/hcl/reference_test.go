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

package hcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceParsing(t *testing.T) {
	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"module", "my-mod"},
			expected: "module.my-mod",
		},
		{
			input:    []string{"aws_s3_bucket", "test"},
			expected: "aws_s3_bucket.test",
		},
		{
			input:    []string{"resource", "aws_s3_bucket", "test"},
			expected: "aws_s3_bucket.test",
		},
		{
			input:    []string{"module", "my-mod"},
			expected: "module.my-mod",
		},
		{
			input:    []string{"data", "aws_iam_policy_document", "s3_policy"},
			expected: "data.aws_iam_policy_document.s3_policy",
		},
		{
			input:    []string{"provider", "aws"},
			expected: "provider.aws",
		},
		{
			input:    []string{"output", "something"},
			expected: "output.something",
		},
	}

	for _, test := range cases {
		t.Run(test.expected, func(t *testing.T) {
			ref, err := newReference(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expected, ref.String())
		})
	}
}
