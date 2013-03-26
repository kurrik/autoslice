// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		err  error
		hs   *AutoSlicer
		args []string
	)
	hs = &AutoSlicer{}
	flag.StringVar(&hs.DstPath, "dst", ".", "Output directory.")
	flag.Parse()
	args = flag.Args()
	if len(args) != 1 {
		fmt.Printf("Must specify target file as first argument\n")
		os.Exit(1)
	}
	hs.SrcPath = args[0]
	if err = hs.Slice(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
