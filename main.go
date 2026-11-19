package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Sample struct {
	Input  string
	Output string
}

var runCmd = map[string]string{
	"go":   "go run main.go",
	"node": "node index.js",
	"bun":  "bun run index.ts",
}

func main() {
	log.SetFlags(0)
	n := flag.Int("n", 0, "Number of problems")
	rt := flag.String("rt", "go", "Runtime environment (go, bun, node)")
	flag.Parse()

	cmd, ok := runCmd[*rt]
	if !ok {
		log.Fatalf("error: unsupported runtime: %s", *rt)
	}

	if *n < 0 || *n > 10 {
		log.Fatal("error: number of problems must be between 1 and 10")
	}

	var targetURL string
	if *n == 0 {
		args := flag.Args()
		if len(args) != 1 {
			log.Fatal("error: exactly one URL is required")
		}

		u, err := url.Parse(args[0])
		if err != nil {
			log.Fatalf("error: invalid URL format: %v", err)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			log.Fatalf("error: unsupported scheme %q (must be http or https)", u.Scheme)
		}

		if u.Hostname() != "atcoder.jp" {
			log.Fatalf("error: invalid host %q (must be atcoder.jp)", u.Hostname())
		}

		targetURL = u.String()
	}

	var samples []Sample
	if targetURL != "" {
		resp, err := http.Get(targetURL)
		if err != nil {
			log.Fatalf("error: failed to fetch url (%v)", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("error: unexpected status code (%d %s)", resp.StatusCode, resp.Status)
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalf("error: failed to parse html (%v)", err)
		}
		var currentInput string
		doc.Find("section").Each(func(i int, s *goquery.Selection) {
			h3 := s.Find("h3").First()
			if h3.Length() == 0 {
				return
			}
			title := h3.Text()
			pre := s.Find("pre").First().Text()
			if strings.HasPrefix(title, "Sample Input") {
				currentInput = pre
			} else if strings.HasPrefix(title, "Sample Output") {
				samples = append(samples, Sample{
					Input:  currentInput,
					Output: pre,
				})
			}
		})
		if len(samples) == 0 {
			log.Fatal("error: no samples found (check contest/number or url)")
		}
	} else {
		for range *n {
			samples = append(samples, Sample{
				Input:  "",
				Output: "",
			})
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	samplesDir := fmt.Sprintf("%s/samples", wd)
	if err = os.RemoveAll(samplesDir); err != nil {
		log.Fatalf("error: failed to remove samples dir (%v)", err)
	}
	var shLines []string
	shLines = append(shLines, "#!/bin/bash", "echo \"=== AtCoder Auto Test ===\"", "")
	for i, v := range samples {
		sampleNum := i + 1
		dirPath := fmt.Sprintf("%s/samples/sample-%d", wd, sampleNum)
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			log.Fatalf("error: failed to create dir %s (%v)", dirPath, err)
		}
		if err = os.WriteFile(fmt.Sprintf("%s/in.txt", dirPath), []byte(v.Input), 0644); err != nil {
			log.Fatalf("error: failed to write in.txt (%v)", err)
		}
		if err = os.WriteFile(fmt.Sprintf("%s/out.txt", dirPath), []byte(v.Output), 0644); err != nil {
			log.Fatalf("error: failed to write out.txt (%v)", err)
		}

		shLines = append(shLines,
			fmt.Sprintf("echo -n \"[Sample %d] \"", sampleNum),
			fmt.Sprintf("if ! %s < ./samples/sample-%d/in.txt > ./samples/sample-%d/actual.txt 2> ./samples/sample-%d/error.txt; then", cmd, sampleNum, sampleNum, sampleNum),
			"  echo -e \"\\e[31m[RE]\\e[0m\"",
			fmt.Sprintf("  cat ./samples/sample-%d/error.txt", sampleNum),
			fmt.Sprintf("elif diff -Z ./samples/sample-%d/out.txt ./samples/sample-%d/actual.txt > /dev/null; then", sampleNum, sampleNum),
			"  echo -e \"\\e[32m[AC]\\e[0m\"",
			"else",
			"  echo -e \"\\e[31m[WA]\\e[0m\"",
			"  echo \"--- Expected Output ---\"",
			fmt.Sprintf("  cat ./samples/sample-%d/out.txt", sampleNum),
			"  echo \"---- Actual Output ----\"",
			fmt.Sprintf("  cat ./samples/sample-%d/actual.txt", sampleNum),
			"  echo \"-----------------------\"",
			"fi",
			"",
		)
	}
	shContent := strings.Join(shLines, "\n")
	shPath := fmt.Sprintf("%s/run_tests.sh", wd)
	if err = os.WriteFile(shPath, []byte(shContent), 0755); err != nil {
		log.Fatalf("error: failed to create run_tests.sh (%v)", err)
	}
	fmt.Printf("\n%d samples(s) saved to %s/samples/\n", len(samples), wd)
	fmt.Printf("Generated %s\n", shPath)
}
