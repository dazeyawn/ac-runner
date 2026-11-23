package main

import (
	"ac-runner/internal"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)

	config, err := internal.ParseConfig()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	samples, err := internal.FetchSamples(config.TargetURL, config.ProblemCount)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	samplesDir := filepath.Join(wd, "samples")
	if err := internal.CleanupDir(samplesDir); err != nil {
		log.Fatalf("error: %v", err)
	}

	var shLines []string
	shLines = append(shLines, "#!/bin/bash", "echo \"=== AtCoder Auto Test ===\"", "")
	for i, v := range samples {
		err = internal.WriteFiles(wd, i+1, v)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	for i := range len(samples) {
		lines := internal.GenScript(config.RunCommand, i+1)
		shLines = append(shLines, lines...)
	}

	shContent := strings.Join(shLines, "\n")
	shPath := filepath.Join(wd, "run_tests.sh")
	if err = os.WriteFile(shPath, []byte(shContent), 0755); err != nil {
		log.Fatalf("error: failed to create run_tests.sh (%v)", err)
	}
	fmt.Printf("\n%d samples(s) saved to %s/samples/\n", len(samples), wd)
	fmt.Printf("Generated %s\n", shPath)
}
