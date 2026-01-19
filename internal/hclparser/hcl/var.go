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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	ctyJson "github.com/zclconf/go-cty/cty/json"
)

func ParseVariable(val interface{}) (cty.Value, error) {
	switch v := val.(type) {
	case string:
		// Try to parse the string as an HCL expression. This will handle expressions
		// passed in via the command line or env variables.
		expr, diags := hclsyntax.ParseExpression([]byte(v), "", hcl.Pos{})
		if !diags.HasErrors() {
			parsedVal, moreDiags := expr.Value(nil)
			if !moreDiags.HasErrors() {
				return parsedVal, nil
			}
		}

		return cty.StringVal(v), nil
	// These cases should only be hit when the input is coming from the config file.
	case int:
		return cty.NumberIntVal(int64(v)), nil
	case float64:
		return cty.NumberFloatVal(v), nil
	case bool:
		return cty.BoolVal(v), nil
	default:
		// Try to parse complex types as JSON
		// This will handle complex variables that have been passed into Infracost
		// via the config file, e.g.:
		// terraform_vars:
		//   my_map:
		//     key1: value1
		//     key2: value2

		// Ensure any maps with non-string keys are converted to string keys
		m := convertToStringKeyMap(v)
		b, err := json.Marshal(m)
		if err != nil {
			return cty.DynamicVal, err
		}

		simple := &ctyJson.SimpleJSONValue{}
		err = simple.UnmarshalJSON(b)
		if err != nil {
			return cty.DynamicVal, err
		}

		return simple.Value, nil
	}
}

func convertToStringKeyMap(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			strKey := fmt.Sprintf("%v", key)
			result[strKey] = convertToStringKeyMap(val)
		}
		return result
	case []interface{}:
		for i, elem := range v {
			v[i] = convertToStringKeyMap(elem)
		}
	}
	return value
}
