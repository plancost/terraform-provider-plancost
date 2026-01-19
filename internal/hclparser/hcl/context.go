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
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/rs/zerolog"
	"github.com/zclconf/go-cty/cty"
)

type Context struct {
	ctx    *hcl.EvalContext
	parent *Context
	logger zerolog.Logger
}

func NewContext(ctx *hcl.EvalContext, parent *Context, logger zerolog.Logger) *Context {
	if ctx.Variables == nil {
		ctx.Variables = make(map[string]cty.Value)
	}

	return &Context{
		ctx:    ctx,
		parent: parent,
		logger: logger,
	}
}

func (c *Context) NewChild() *Context {
	return NewContext(c.ctx.NewChild(), c, c.logger)
}

func (c *Context) Parent() *Context {
	return c.parent
}

func (c *Context) Inner() *hcl.EvalContext {
	return c.ctx
}

func (c *Context) Root() *Context {
	root := c
	for {
		if root.Parent() == nil { //nolint:staticcheck
			break
		}
		root = root.Parent()
	}
	return root
}

func (c *Context) Get(parts ...string) (val cty.Value) {
	defer func() {
		if val == cty.NilVal && c.Parent() != nil {
			val = c.Parent().Get(parts...)
		}
	}()

	if len(parts) == 0 {
		return cty.NilVal
	}

	src := c.ctx.Variables
	for i, part := range parts {
		if i == len(parts)-1 {
			return src[part]
		}
		nextPart := src[part]
		if nextPart == cty.NilVal {
			return cty.NilVal
		}
		src = nextPart.AsValueMap()
	}

	return cty.NilVal
}

func (c *Context) SetByDot(val cty.Value, path string) {
	c.Set(val, strings.Split(path, ".")...)
}

func (c *Context) Set(val cty.Value, parts ...string) {
	if len(parts) == 0 {
		return
	}

	v := mergeVars(c.ctx.Variables[parts[0]], parts[1:], val)
	c.ctx.Variables[parts[0]] = v
}

func isValidCtyObject(src cty.Value) bool {
	return src.IsKnown() && src.Type().IsObjectType() && !src.IsNull() && src.LengthInt() > 0
}

func mergeVars(src cty.Value, parts []string, value cty.Value) cty.Value {
	if len(parts) == 0 {
		if isValidCtyObject(value) && isValidCtyObject(src) {
			return mergeObjects(src, value)
		}

		return value
	}

	var data map[string]cty.Value
	if src.Type().IsObjectType() && !src.IsNull() && src.LengthInt() > 0 {
		data = src.AsValueMap()
		tmp, ok := data[parts[0]] // nolint:gosec // ignore "G602: slice index out of range" since we already check len(parts) == 0
		if !ok {
			src = cty.ObjectVal(make(map[string]cty.Value))
		} else {
			src = tmp
		}
	} else {
		data = make(map[string]cty.Value)
	}

	data[parts[0]] = mergeVars(src, parts[1:], value) // nolint:gosec // ignore "G602: slice index out of range" since we already check len(parts) == 0

	return cty.ObjectVal(data)
}

// mergeObjects merges two cty.Value objects by recursively combining their
// key-value pairs. When there are conflicting keys, the value from object `b`
// takes precedence over object `a`, unless both values are valid cty objects, in
// which case they are recursively merged.
func mergeObjects(a cty.Value, b cty.Value) cty.Value {
	output := a.AsValueMap()
	for key, val := range b.AsValueMap() {
		old, exists := output[key]

		if exists && isValidCtyObject(val) && isValidCtyObject(old) {
			output[key] = mergeObjects(old, val)
			continue
		}

		output[key] = val
	}
	return cty.ObjectVal(output)
}
