// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/pprof/profile"
)

func main() {
	outputFile := flag.String("o", "", "output file path")
	flag.Parse()

	// Parse the command-line arguments
	flag.Parse()

	// Get the positional arguments
	files := flag.Args()

	// Check if the output file flag is provided
	if *outputFile == "" {
		log.Panicln("Output file path is required")
	}

	if len(files) == 0 {
		log.Panicln("Give profiles files as arguments")
	}

	var profiles []*profile.Profile
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			log.Panicf("Cannot open profile file at %q: %v", fname, err)
		}
		p, err := profile.Parse(f)
		if err != nil {
			log.Panicf("Cannot parse profile at %q: %v", fname, err)
		}
		profiles = append(profiles, p)
	}

	merged, err := profile.Merge(profiles)
	if err != nil {
		log.Fatalf("Cannot merge profiles: %v", err)
	}

	out, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Panicf("Cannot open output to write: %v", err)
	}

	if err := merged.Write(out); err != nil {
		log.Panicf("Cannot write merged profile to file: %v", err)
	}

	if err := out.Close(); err != nil {
		log.Panicf("Error when closing the output file: %v", err)
	}
}
