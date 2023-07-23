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

package main

import (
	"flag"
	"github.com/chenmingyong0423/gkit/stringx"
	"github.com/chenmingyong0423/go-optioner/options"
	"log"
	"os"
)

var (
	structTypeNameArg = flag.String("type", "", "Struct type name of the functional options struct.")
	outputArg         = flag.String("output", "", "Output file name, default: srcDir/opt_<struct type>_gen.go")
	g                 = options.NewGenerator()
)

func usage() {
	log.Printf("go-optioner is a tool for generating functional options pattern.\n")
	log.Printf("Usage: \n")
	log.Printf("\t go-optioner [flags]\n")
	log.Printf("Flags:\n")
	log.Printf("\t -type <struct name>\n")
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(*structTypeNameArg) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	g.StructInfo.StructName = *structTypeNameArg
	g.StructInfo.NewStructName = stringx.BigCamelToSmallCamel(*structTypeNameArg)
	g.SetOutPath(outputArg)

	g.GeneratingOptions()
	if !g.Found {
		log.Printf("Target \"[%s]\" is not be found\n", g.StructInfo.StructName)
		os.Exit(1)
	}

	g.GenerateCodeByTemplate()
	g.OutputToFile()
}
