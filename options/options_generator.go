// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package options

import (
	"bytes"
	"fmt"
	"github.com/chenmingyong0423/gkit/stringx"
	"github.com/chenmingyong0423/go-optioner/templates"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"
)

type Generator struct {
	StructInfo *StructInfo
	outPath    string

	buf   bytes.Buffer
	Found bool
}

func NewGenerator() *Generator {
	return &Generator{
		StructInfo: &StructInfo{
			Fields:         make([]FieldInfo, 0),
			OptionalFields: make([]FieldInfo, 0),
		},
	}
}

type FieldInfo struct {
	Name string
	Type string
}

type StructInfo struct {
	PackageName    string
	StructName     string
	NewStructName  string
	Fields         []FieldInfo
	OptionalFields []FieldInfo
}

func (g *Generator) GeneratingOptions() {
	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatalf("Processsing directory failed: %s", err.Error())
	}
	for _, file := range pkg.GoFiles {
		if found := g.parseStruct(file); found {
			g.Found = found
			break
		}
	}
}

func (g *Generator) parseStruct(fileName string) bool {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, fileName, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	g.StructInfo.PackageName = file.Name.Name

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if typeSpec.Name.String() != g.StructInfo.StructName {
				continue
			}
			if structDecl, ok := typeSpec.Type.(*ast.StructType); ok {
				log.Printf("Generating Struct \"%s\" \n", g.StructInfo.StructName)
				for _, field := range structDecl.Fields.List {
					if len(field.Names) == 0 {
						continue
					}

					optionIgnore := false

					fieldName := field.Names[0].Name
					fieldType := field.Type.(*ast.Ident).Name
					if field.Tag != nil {
						tags := strings.Replace(field.Tag.Value, "`", "", -1)
						tag := reflect.StructTag(tags).Get("opt")
						if tag == "-" {
							g.StructInfo.Fields = append(g.StructInfo.Fields, FieldInfo{
								Name: fieldName,
								Type: fieldType,
							})
							optionIgnore = true
						}
					}
					if !optionIgnore {
						g.StructInfo.OptionalFields = append(g.StructInfo.OptionalFields, FieldInfo{
							Name: fieldName,
							Type: fieldType,
						})
					}
					log.Printf("Generating Struct Field \"%s\" of type \"%s\"\n", fieldName, fieldType)
				}
				return true
			} else {
				log.Fatal(fmt.Sprintf("Target[%s] type is not a struct", g.StructInfo.StructName))
			}
		}
	}
	return false
}

func (g *Generator) GenerateCodeByTemplate() {
	tmpl, err := template.New("options").Funcs(template.FuncMap{"bigCamelToSmallCamel": stringx.BigCamelToSmallCamel}).Parse(templates.OptionsTemplateCode)
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		os.Exit(1)
	}

	err = tmpl.Execute(&g.buf, g.StructInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Generator) OutputToFile() {
	src := g.forMart()
	err := os.WriteFile(g.outPath, src, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generating Functional Options Code Successful.\nOut: %s\n", g.outPath)
}

func (g *Generator) forMart() []byte {
	source, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return source
}

func (g *Generator) SetOutPath(outPath *string) {
	if len(*outPath) > 0 {
		g.outPath = *outPath
	} else {
		g.outPath = fmt.Sprintf("opt_%s_gen.go", stringx.CamelToSnake(g.StructInfo.StructName))
	}
}
