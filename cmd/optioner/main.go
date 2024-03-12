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
	"fmt"
	"github.com/chenmingyong0423/gkit/stringx"
	"github.com/chenmingyong0423/go-optioner/options"
	"log"
	"os"
)

type ModeValue struct {
	value       string
	validValues []string
}

func (m *ModeValue) String() string {
	return m.value
}

func (m *ModeValue) Set(s string) error {
	for _, v := range m.validValues {
		if s == v {
			m.value = s
			return nil
		}
	}
	return fmt.Errorf("invalid value %q for mode, valid values are: %v", s, m.validValues)
}

var (
	outputMode = ModeValue{
		value:       "write",
		validValues: []string{"write", "append"},
	}
	structTypeName = flag.String("type", "", "Struct type name of the functional options struct.")
	output         = flag.String("output", "", "Output file name, default: srcDir/opt_<struct type>_gen.go")
	g              = options.NewGenerator()
)

func usage() {
	fmt.Fprintf(os.Stderr, "optioner is a tool for generating functional options pattern.\n")
	fmt.Fprintf(os.Stderr, "Usage: \n")
	fmt.Fprintf(os.Stderr, "\t optioner [flags]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	fmt.Fprintf(os.Stderr, "\t -type <struct name>\n")
	fmt.Fprintf(os.Stderr, "\t -output <output path>, default: srcDir/opt_xxx_gen.go\n")
	fmt.Fprintf(os.Stderr, "\t -mode <the file writing mode>, default: write\n")
	fmt.Fprintf(os.Stderr, "\t there are two available modes:\n")
	fmt.Fprintf(os.Stderr, "\t\t - write(Write/Overwrite): Overwrites or creates a new file.\n")
	fmt.Fprintf(os.Stderr, "\t\t - append (Append): Adds to the end of the file.\n")

}

func main() {
	flag.Var(&outputMode, "mode", "The file writing mode, default: write")
	flag.Usage = usage
	flag.Parse()
	if len(*structTypeName) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	g.StructInfo.StructName = *structTypeName
	g.StructInfo.NewStructName = stringx.BigCamelToSmallCamel(*structTypeName)
	g.SetOutPath(output)
	g.SetMod(outputMode.value)

	g.GeneratingOptions()
	if !g.Found {
		log.Printf("Target \"[%s]\" is not be found\n", g.StructInfo.StructName)
		os.Exit(1)
	}
	g.GenerateCodeByTemplate()
	g.OutputToFile()
}
