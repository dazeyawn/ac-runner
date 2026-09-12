package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Sample struct {
	Input  string
	Output string
}

func main() {
	log.SetFlags(0)
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal(`usage: command <contest> <problem>
   or: command -u <url>
   or: command -n <number>`)
	}

	var url string
	switch args[0] {
	case "-u":
		if !strings.HasPrefix(args[1], "https://atcoder.jp") {
			log.Fatal("error: URL must start with https://atcoder.jp")
		}
		url = args[1]
	case "-n":
		url = ""
	default:
		contest := args[0]
		problem := args[1]
		if len(problem) != 1 {
			log.Fatalf("error: problem must be single character (got %s)", problem)
		}
		url = fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s_%s", contest, contest, problem)
	}

	var samples []Sample
	if url != "" {
		resp, err := http.Get(url)
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
		n, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("error: -n requires a number (got: %s)", args[1])
		}
		if n < 1 {
			log.Fatalf("error: number of samples must be at least 1 (got: %d)", n)
		}
		for range n {
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
			fmt.Sprintf("go run main.go < ./samples/sample-%d/in.txt > ./samples/sample-%d/actual.txt 2>&1", sampleNum, sampleNum),
			fmt.Sprintf("diff -Z ./samples/sample-%d/out.txt ./samples/sample-%d/actual.txt > /dev/null", sampleNum, sampleNum),
			"if [ $? -eq 0 ]; then",
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
