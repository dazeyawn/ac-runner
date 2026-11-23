package internal

type Config struct {
	RunCommand   string
	ProblemCount int
	TargetURL    string
}

var runCmd = map[string]string{
	"go":   "go run main.go",
	"node": "node index.js",
	"bun":  "bun run index.ts",
}

type Sample struct {
	Input  string
	Output string
}
