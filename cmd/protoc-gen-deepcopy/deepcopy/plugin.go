// Copyright 2019 Istio Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package deepcopy

import (
	"path"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func init() {
	generator.RegisterPlugin(NewPlugin())
}

// FileNameSuffix is the suffix added to files generated by deepcopy
const FileNameSuffix = "_deepcopy.gen.go"
const kubeTypeGenTag = "+kubetype-gen"

// Plugin is a protoc-gen-gogo plugin that creates DeepCopyInto() functions for
// annotated protobuf types.
type Plugin struct {
	*generator.Generator
	generator.PluginImports
	filesWritten map[string]interface{}
}

// NewPlugin returns a new instance of the Plugin
func NewPlugin() *Plugin {
	return &Plugin{
		filesWritten: map[string]interface{}{},
	}
}

// Name returns the name of this plugin
func (p *Plugin) Name() string {
	return "deepcopy"
}

// Init initializes our plugin with the active generator
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
}

// Generate our content
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	// imported packages
	// XXX: proto is imported by default
	//protoPkg := p.NewImport("github.com/gogo/protobuf/proto")

	wroteMarshalers := false

	var leadingDetachedCommentsMap map[string][]string

messageLoop:
	for _, message := range file.Messages() {
		// check to make sure something was generated for this type
		if !gogoproto.HasTypeDecl(file.FileDescriptorProto, message.DescriptorProto) {
			continue
		}

		comments := p.Comments(message.Path())
		if !strings.Contains(comments, kubeTypeGenTag) {
			// check leading detached comments
			if leadingDetachedCommentsMap == nil {
				leadingDetachedCommentsMap = buildLeadingDetachedCommentsMap(file)
			}
			var tagged bool
			for _, detachedComment := range leadingDetachedCommentsMap[message.Path()] {
				if strings.Contains(detachedComment, kubeTypeGenTag) {
					tagged = true
					break
				}
			}
			if !tagged {
				continue messageLoop
			}
		}

		typeName := generator.CamelCaseSlice(message.TypeName())

		// Generate DeepCopyInto() method for this type
		p.P(`// DeepCopyInto supports using `, typeName, ` within kubernetes types, where deepcopy-gen is used.`)
		p.P(`func (in *`, typeName, `) DeepCopyInto(out *`, typeName, `) {`)
		p.In()
		p.P(`p := proto.Clone(in).(*`, typeName, `)`)
		p.P(`*out = *p`)
		p.Out()
		p.P(`}`)

		// Generate DeepCopy() method for this type
		p.P(`// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new `, typeName, `. Required by controller-gen.`)
		p.P(`func (in *`, typeName, `) DeepCopy() *`, typeName, ` {`)
		p.In()
		p.P(`if in == nil { return nil }`)
		p.P(`out := new(`, typeName, `)`)
		p.P(`in.DeepCopyInto(out)`)
		p.P(`return out`)
		p.Out()
		p.P(`}`)

		wroteMarshalers = true
	}

	if !wroteMarshalers {
		return
	}

	// store this file away
	p.addFile(file)
}

func (p *Plugin) addFile(file *generator.FileDescriptor) {
	name := file.GetName()
	importPath := ""
	// the relevant bits of FileDescriptor.goPackageOption(), if only it were exported
	opt := file.GetOptions().GetGoPackage()
	if opt != "" {
		if sc := strings.Index(opt, ";"); sc >= 0 {
			// A semicolon-delimited suffix delimits the import path and package name.
			importPath = opt[:sc]
		} else if strings.LastIndex(opt, "/") > 0 {
			// The presence of a slash implies there's an import path.
			importPath = opt
		}
	}
	// strip the extension
	name = name[:len(name)-len(path.Ext(name))]
	if importPath != "" {
		name = path.Join(importPath, path.Base(name))
	}
	p.filesWritten[name+FileNameSuffix] = struct{}{}
}

// FilesWritten returns a list of the names of files for which output was generated
func (p *Plugin) FilesWritten() map[string]interface{} {
	return p.filesWritten
}

// buildLeadingDetachedCommentsMap returns a map of source code location path => leading detached comments
// the location path is converted to a string (source is []int32).  this is used
// above to quickly retrieve leading detached comments for messages (message.Path(), which is a stringified
// version of the source code location for a message).
func buildLeadingDetachedCommentsMap(file *generator.FileDescriptor) map[string][]string {
	retVal := make(map[string][]string)
	for _, location := range file.GetSourceCodeInfo().GetLocation() {
		ldc := location.GetLeadingDetachedComments()
		if len(ldc) > 0 {
			intpath := location.GetPath()
			strpath := make([]string, 0, len(intpath))
			for _, intelement := range intpath {
				strpath = append(strpath, strconv.FormatInt(int64(intelement), 10))
			}
			retVal[strings.Join(strpath, ",")] = ldc
		}
	}
	return retVal
}
