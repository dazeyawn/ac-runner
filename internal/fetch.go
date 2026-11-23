package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchSamples(targetURL string, problemCount int) ([]Sample, error) {
	var samples []Sample

	if targetURL == "" {
		for range problemCount {
			samples = append(samples, Sample{
				Input:  "",
				Output: "",
			})
		}
		return samples, nil
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		return samples, fmt.Errorf("failed to fetch url (%v)", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return samples, fmt.Errorf("unexpected status code (%d %s)", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return samples, fmt.Errorf("failed to parse html (%v)", err)
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
		return samples, fmt.Errorf("no samples found (check contest/number or url)")
	}

	return samples, nil
}
