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

package funcs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	componentsFuncs "github.com/turbot/terraform-components/lang/funcs"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
)

func MakeFileFunc(baseDir string, encBase64 bool) function.Function {
	ff := componentsFuncs.MakeFileFunc(baseDir, encBase64)

	return function.New(&function.Spec{
		Params: ff.Params(),
		Type:   function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			c, err := ff.Call(args)

			// if we get an error calling the underlying file function this is likely because the path
			// argument has been transformed at some point because of infracost mocking/state fallbacks
			// so we return a blank string instead of failing the evaluation. This is safer than returning
			// an error as we in complex expression cases we can actually return a partial value instead
			// of a unknown value which will cause subsequent evaluations to fail.
			if err != nil {
				logging.Logger.Debug().Msgf("error calling file func: %s returning a blank string for filesytem func", err)
				return cty.StringVal(""), nil
			}

			return c, nil
		},
	})
}

// MakeTemplateFileFunc constructs a function that takes a file path and
// an arbitrary object of named values and attempts to render the referenced
// file as a template using HCL template syntax.
//
// The template itself may recursively call other functions so a callback
// must be provided to get access to those functions. The template cannot,
// however, access any variables defined in the scope: it is restricted only to
// those variables provided in the second function argument, to ensure that all
// dependencies on other graph nodes can be seen before executing this function.
//
// As a special exception, a referenced template file may not recursively call
// the templatefile function, since that would risk the same file being
// included into itself indefinitely.
func MakeTemplateFileFunc(baseDir string, funcsCb func() map[string]function.Function) function.Function {
	ff := componentsFuncs.MakeTemplateFileFunc(baseDir, funcsCb)
	return function.New(&function.Spec{
		Params: ff.Params(),
		Type: func(args []cty.Value) (cty.Type, error) {
			if !args[0].IsKnown() || !args[1].IsKnown() {
				return cty.DynamicPseudoType, nil
			}

			return ff.ReturnTypeForValues(args)
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			return ff.Call(args)
		},
	})

}

// MakeFileExistsFunc constructs a function that takes a path
// and determines whether a file exists at that path
func MakeFileExistsFunc(baseDir string) function.Function {
	ff := componentsFuncs.MakeFileExistsFunc(baseDir)
	return function.New(&function.Spec{
		Params: ff.Params(),
		Type:   function.StaticReturnType(cty.Bool),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			return ff.Call(args)
		},
	})
}

// MakeFileSetFunc constructs a function that takes a glob pattern
// and enumerates a file set from that pattern
func MakeFileSetFunc(baseDir string) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: cty.String,
			},
			{
				Name: "pattern",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Set(cty.String)),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			path := args[0].AsString()
			pattern := args[1].AsString()

			if !filepath.IsAbs(path) {
				path = filepath.Join(baseDir, path)
			}

			pattern = filepath.Join(path, pattern)
			matches, err := doublestar.Glob(pattern)
			if err != nil {
				return cty.UnknownVal(cty.Set(cty.String)), fmt.Errorf("failed to glob pattern (%s): %s", pattern, err)
			}

			var matchVals []cty.Value
			for _, match := range matches {
				fi, err := os.Stat(match)

				if err != nil {
					return cty.UnknownVal(cty.Set(cty.String)), fmt.Errorf("failed to stat (%s): %s", match, err)
				}

				if !fi.Mode().IsRegular() {
					continue
				}

				match, err = filepath.Rel(path, match)
				if err != nil {
					return cty.UnknownVal(cty.Set(cty.String)), fmt.Errorf("failed to trim path of match (%s): %s", match, err)
				}

				match = filepath.ToSlash(match)
				matchVals = append(matchVals, cty.StringVal(match))
			}

			if len(matchVals) == 0 {
				return cty.SetValEmpty(cty.String), nil
			}

			return cty.SetVal(matchVals), nil
		},
	})
}

// File reads the contents of the file at the given path.
//
// The file must contain valid UTF-8 bytes, or this function will return an error.
//
// The underlying function implementation works relative to a particular base
// directory, so this wrapper takes a base directory string and uses it to
// construct the underlying function before calling it.
func File(baseDir string, path cty.Value) (cty.Value, error) {
	fn := MakeFileFunc(baseDir, false)
	return fn.Call([]cty.Value{path})
}

// FileExists determines whether a file exists at the given path.
//
// The underlying function implementation works relative to a particular base
// directory, so this wrapper takes a base directory string and uses it to
// construct the underlying function before calling it.
func FileExists(baseDir string, path cty.Value) (cty.Value, error) {
	fn := MakeFileExistsFunc(baseDir)
	return fn.Call([]cty.Value{path})
}

// FileSet enumerates a set of files given a glob pattern
//
// The underlying function implementation works relative to a particular base
// directory, so this wrapper takes a base directory string and uses it to
// construct the underlying function before calling it.
func FileSet(baseDir string, path, pattern cty.Value) (cty.Value, error) {
	fn := MakeFileSetFunc(baseDir)
	return fn.Call([]cty.Value{path, pattern})
}

// FileBase64 reads the contents of the file at the given path.
//
// The bytes from the file are encoded as base64 before returning.
//
// The underlying function implementation works relative to a particular base
// directory, so this wrapper takes a base directory string and uses it to
// construct the underlying function before calling it.
func FileBase64(baseDir string, path cty.Value) (cty.Value, error) {
	fn := MakeFileFunc(baseDir, true)
	return fn.Call([]cty.Value{path})
}
