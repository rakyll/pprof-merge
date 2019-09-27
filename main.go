package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/pprof/profile"
)

var (
	output string
)

func main() {
	flag.StringVar(&output, "o", "merged.data", "")
	flag.Parse()

	files := os.Args[1:]
	if len(files) == 0 {
		log.Fatal("Give profiles files as arguments")
	}

	var profiles []*profile.Profile
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatalf("Cannot open profile file at %q: %v", fname, err)
		}
		p, err := profile.Parse(f)
		if err != nil {
			log.Fatalf("Cannot parse profile at %q: %v", fname, err)
		}
		profiles = append(profiles, p)
	}

	merged, err := profile.Merge(profiles)
	if err != nil {
		log.Fatalf("Cannot merge profiles: %v", err)
	}

	out, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Cannot open output to write: %v", err)
	}

	if err := merged.Write(out); err != nil {
		log.Fatalf("Cannot write merged profile to file: %v", err)
	}
}
