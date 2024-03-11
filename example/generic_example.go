// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package example

import "context"

// 用于添加 "context" 导包信息，以便在生成代码时判断是否成功去除该无用包信息
var _ context.Context

//go:generate go run ../cmd/optioner/main.go -type GenericExample
type GenericExample[T any, U comparable, V ~int] struct {
	A T `opt:"-"`
	B U
	C V
	D string
}
