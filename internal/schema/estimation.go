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

package schema

import "context"

// EstimateFunc queries cloud providers to gather past usage information, then
// projects future usage based on the results.
type EstimateFunc func(context.Context, map[string]interface{}) error

// Remediater allows correction of cloud configuration issues
// so that future runs of Infracost will provide more accurate results.
type Remediater interface {
	// Describe provides an English description of the remediation action X that
	// would fit into a sentence "May we X?" (e.g. "enable bucket metrics").
	// The description can be used to prompt the user before taking action.
	Describe() string

	// Remediate attempts to fix a problem in the cloud that prevents estimation,
	// e.g. by enabling metrics collection on certain resources.
	Remediate() error
}
