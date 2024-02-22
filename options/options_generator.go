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
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/chenmingyong0423/gkit/stringx"
	"github.com/chenmingyong0423/go-optioner/templates"
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
	GenericParams  []FieldInfo

	Imports []string
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
				if typeSpec.TypeParams != nil {
					log.Println("This is a struct which contains generic type:", typeSpec.Name)
					for _, param := range typeSpec.TypeParams.List {
						for _, name := range param.Names {
							typ := g.getTypeName(param.Type)
							g.StructInfo.GenericParams = append(g.StructInfo.GenericParams, FieldInfo{
								Name: name.Name,
								Type: typ,
							})
							log.Printf("Generic parameter: %s %s\n", name.Name, typ)
						}
					}
				}
				for _, field := range structDecl.Fields.List {
					fieldName := ""
					if len(field.Names) == 0 {
						if ident, ok := field.Type.(*ast.Ident); ok { // combined struct
							fieldName = ident.Name
						} else {
							continue
						}
					} else {
						fieldName = field.Names[0].Name
					}
					optionIgnore := false

					fieldType := g.getTypeName(field.Type)
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
				// 收集 package 信息
				g.CollectImports(file)
				return true
			} else {
				log.Fatal(fmt.Sprintf("Target[%s] type is not a struct", g.StructInfo.StructName))
			}
		}
	}
	return false
}

func (g *Generator) GenerateCodeByTemplate() {
	tmpl, err := template.New("options").Funcs(template.FuncMap{"bigCamelToSmallCamel": stringx.BigCamelToSmallCamel, "capitalizeFirstLetter": stringx.CapitalizeFirstLetter}).Parse(templates.OptionsTemplateCode)
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
	fmt.Printf("Generating Functional Options Code Successfully.\nOut: %s\n", g.outPath)
}

func (g *Generator) forMart() []byte {
	source, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return source
}

func (g *Generator) SetOutPath(outPath *string) {
	fileName := fmt.Sprintf("opt_%s_gen.go", stringx.CamelToSnake(g.StructInfo.StructName))
	if len(*outPath) > 0 {
		g.outPath = *outPath
	} else {
		g.outPath = fileName
	}
}

func (g *Generator) getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", g.getTypeName(t.X), t.Sel.Name)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + g.getTypeName(t.Elt)
		}
		if basicLit, ok := t.Len.(*ast.BasicLit); ok && basicLit.Kind == token.INT {
			return "[" + basicLit.Value + "]" + g.getTypeName(t.Elt)
		} else {
			log.Fatalf("Array len error: %T", t)
			return ""
		}
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", g.getTypeName(t.Key), g.getTypeName(t.Value))
	case *ast.StarExpr:
		return "*" + g.getTypeName(t.X)
	//case *ast.InterfaceType:
	//	return "" // ignore
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return g.parseFuncType(t)
	case *ast.ChanType:
		return "chan " + g.getTypeName(t.Value)
	case *ast.UnaryExpr:
		return "~" + g.getTypeName(t.X)
	default:
		log.Fatalf("Unsupported type for field: %T", t)
		return ""
	}
}

func (g *Generator) parseFuncType(f *ast.FuncType) string {
	var params, results []string
	if f.Params != nil {
		for _, field := range f.Params.List {
			paramType := g.getTypeName(field.Type)
			for _, name := range field.Names {
				params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
			}
		}
	}

	if f.Results != nil {
		for _, field := range f.Results.List {
			resultType := g.getTypeName(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					results = append(results, fmt.Sprintf("%s %s", name.Name, resultType))
				}
			} else {
				results = append(results, resultType)
			}
		}
	}

	if len(results) == 1 {
		return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), results[0])
	}
	return fmt.Sprintf("func(%s) (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}

func (g *Generator) CollectImports(file *ast.File) {
	for _, imp := range file.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			log.Fatalf("Failed to unquote import path: %v", err)
		}
		g.StructInfo.Imports = append(g.StructInfo.Imports, path)
	}
}
