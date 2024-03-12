// Generated by [optioner] command-line tool; DO NOT EDIT
// If you have any questions, please create issues and submit contributions at:
// https://github.com/chenmingyong0423/go-optioner

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

type _ struct {
	_ context.Context
}

type GenericExampleOption[T any, U comparable, V ~int] func(*GenericExample[T, U, V])

func NewGenericExample[T any, U comparable, V ~int](a T, opts ...GenericExampleOption[T, U, V]) *GenericExample[T, U, V] {
	genericExample := &GenericExample[T, U, V]{
		A: a,
	}

	for _, opt := range opts {
		opt(genericExample)
	}

	return genericExample
}

func WithB[T any, U comparable, V ~int](b U) GenericExampleOption[T, U, V] {
	return func(genericExample *GenericExample[T, U, V]) {
		genericExample.B = b
	}
}

func WithC[T any, U comparable, V ~int](c V) GenericExampleOption[T, U, V] {
	return func(genericExample *GenericExample[T, U, V]) {
		genericExample.C = c
	}
}

func WithD[T any, U comparable, V ~int](d string) GenericExampleOption[T, U, V] {
	return func(genericExample *GenericExample[T, U, V]) {
		genericExample.D = d
	}
}