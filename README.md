# ac-runner

A Go-based CLI tool that fetches sample test cases from AtCoder problem pages and generates a `run_tests.sh` script to run tests using your preferred runtime.

## Features

- Automatically fetches sample inputs and outputs from `atcoder.jp` problem URLs.
- Creates empty sample directories (`samples/sample-1`, etc.) with the `-n` option.
- Generates a `run_tests.sh` script with runtime-specific commands using the `-rt` option (`go`, `bun`, or `node`).
- Displays test results (`[AC]`, `[WA]`, `[RE]`) and output differences in the terminal.

## Usage

### 1. Fetch samples and generate the test script

```bash
go run . -rt go https://atcoder.jp/contests/abc001/tasks/abc001_1
```

You can also specify a different runtime. The default is `go`.

```bash
go run . -rt bun https://atcoder.jp/contests/abc001/tasks/abc001_1
go run . -rt node https://atcoder.jp/contests/abc001/tasks/abc001_1
```

### 2. Create empty sample directories

Use the `-n` option to create a specified number of empty sample directories (1–10).

```bash
go run . -n 3 -rt go
```

### 3. Run tests

```bash
./run_tests.sh
```

## Generated Structure

```text
.
├── run_tests.sh
└── samples/
    ├── sample-1/
    │   ├── in.txt
    │   └── out.txt
    └── sample-2/
        ├── in.txt
        └── out.txt
```

**Note:** Running `./run_tests.sh` also generates `actual.txt` and `error.txt` files.

## Supported Runtimes

| `-rt` | Command | Target File |
|-------|---------|-------------|
| `go` | `go run main.go` | `main.go` |
| `node` | `node index.js` | `index.js` |
| `bun` | `bun run index.ts` | `index.ts` |
