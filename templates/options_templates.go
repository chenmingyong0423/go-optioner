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

package templates

const OptionsTemplateCode = `
// Generated by optioner -type {{ .StructName }}; DO NOT EDIT
// If you have any questions, please create issues and submit contributions at:
// https://github.com/chenmingyong0423/go-optioner

package {{ .PackageName }}

{{ if .Imports }}
import (
    {{ range .Imports }}
    "{{ . }}"
    {{ end }}
)
{{ end }}

type {{ .StructName }}Option func(*{{ .StructName }})

func New{{ .StructName }}({{ range $index, $field := .Fields }}{{ $field.Name | bigCamelToSmallCamel }} {{ $field.Type }},{{ end }} opts ...{{ .StructName }}Option) *{{ .StructName }} {
	{{ .NewStructName }} := &{{ .StructName }}{
		{{ range $index, $field := .Fields }}{{ $field.Name }}: {{ $field.Name | bigCamelToSmallCamel }},
		{{ end }}
	}

	for _, opt := range opts {
		opt({{ .NewStructName }})
	}

	return {{ .NewStructName }}
}

{{ if .OptionalFields }}
{{ range $field := .OptionalFields }}
func With{{ $field.Name | capitalizeFirstLetter }}({{ $field.Name | bigCamelToSmallCamel }} {{ $field.Type }}) {{ $.StructName }}Option {
	return func({{ $.NewStructName }} *{{ $.StructName }}) {
		{{ $.NewStructName }}.{{ $field.Name }} = {{ $field.Name | bigCamelToSmallCamel }}
	}
}
{{ end }}
{{ end }}
`
