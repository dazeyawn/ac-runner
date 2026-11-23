package internal

import "fmt"

func GenScript(runCommand string, sampleNum int) []string {
	return []string{
		fmt.Sprintf("echo -n \"[Sample %d] \"", sampleNum),
		fmt.Sprintf("if ! %s < ./samples/sample-%d/in.txt > ./samples/sample-%d/actual.txt 2> ./samples/sample-%d/error.txt; then", runCommand, sampleNum, sampleNum, sampleNum),
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
	}
}
