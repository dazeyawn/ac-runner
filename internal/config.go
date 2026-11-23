package internal

import (
	"flag"
	"fmt"
	"net/url"
)

func ParseConfig() (Config, error) {
	var config Config

	n := flag.Int("n", 0, "Number of problems")
	rt := flag.String("rt", "go", "Runtime environment (go, bun, node)")
	flag.Parse()

	cmd, ok := runCmd[*rt]
	if !ok {
		return config, fmt.Errorf("unsupported runtime: %s", *rt)
	}
	config.RunCommand = cmd

	if *n < 0 || *n > 10 {
		return config, fmt.Errorf("number of problems must be 0 (url mode) or between 1 and 10")
	}
	config.ProblemCount = *n

	if *n == 0 {
		args := flag.Args()
		if len(args) != 1 {
			return config, fmt.Errorf("exactly one URL is required")
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return config, fmt.Errorf("invalid URL format: %v", err)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return config, fmt.Errorf("unsupported scheme %q (must be http or https)", u.Scheme)
		}

		if u.Hostname() != "atcoder.jp" {
			return config, fmt.Errorf("invalid host %q (must be atcoder.jp)", u.Hostname())
		}

		config.TargetURL = u.String()
	}

	return config, nil
}
