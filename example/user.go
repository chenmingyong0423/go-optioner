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

package example

import (
	"context"
	"github.com/chenmingyong0423/go-optioner/example/third_party"
)

// 用于添加 "context" 导包信息，以便在生成代码时判断是否成功去除该无用包信息
var _ context.Context

type Embedded struct{}

type Embedded2 struct{}

type Embedded3 struct{}

type Embedded4 struct{}

type Embedded5 struct{}

type Embedded6 struct{}

type Embedded7 struct{}

type Embedded8 struct{}

//go:generate go run ../cmd/optioner/main.go -type GenericExample
type GenericExample[T any, U comparable, V ~int] struct {
	A T `opt:"-"`
	B U
	C V
	D string
}

//go:generate go run ../cmd/optioner/main.go -type User
type User struct {
	Embedded   `opt:"-"`
	*Embedded2 `opt:"-"`
	E3         Embedded3  `opt:"-"`
	E4         *Embedded4 `opt:"-"`
	Embedded5
	*Embedded6
	E7                 Embedded7
	E8                 *Embedded8
	Username           string
	Email              string
	Address            // combined struct
	ArrayField         [4]int
	SliceField         []int
	ThirdPartyField    third_party.ThirdParty
	MapField           map[string]int
	PtrField           *int
	EmptyStructFiled   struct{}
	SimpleFuncField    func()
	ComplexFuncField   func(a int)
	ComplexFuncFieldV2 func() int
	ComplexFuncFieldV3 func(a int) int
	ComplexFuncFieldV4 func(a int) (int, error)
	ChanField          chan int
	error              // interface
}

type Address struct {
	Street string
	City   string
}
