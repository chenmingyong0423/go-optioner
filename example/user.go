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
	"github.com/chenmingyong0423/go-optioner/example/third_party"
)

type User struct {
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
